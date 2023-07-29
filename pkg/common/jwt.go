package common

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/sync/singleflight"
	"time"
)

var concurrency_control = &singleflight.Group{}

type JWT struct {
	Issuer     string
	Key        []byte
	Timeout    time.Duration
	BufferTime time.Duration
}

var (
	TokenTimeout   = errors.New("Token is timeout")
	TokenNotActive = errors.New("Token not active")
	TokenInvalid   = errors.New("Token invalid")
)

func NewJWT(key string, timeout int, Issuer string, bufferTime int) *JWT {
	j := &JWT{Key: []byte(key), Timeout: time.Duration(timeout) * time.Second, Issuer: Issuer}
	j.BufferTime = j.Timeout - (time.Duration(bufferTime) * time.Second)
	return j
}

type BaseClaims struct {
	UserID   uint
	Account  string
	UUID     string
	ClientId string
}
type CustomClaims struct {
	BaseClaims
	BufferTime time.Time
	jwt.RegisteredClaims
}

func (j *JWT) CreateClaims(base BaseClaims) CustomClaims {
	claims := CustomClaims{
		BaseClaims: base,
		BufferTime: time.Now().Add(j.BufferTime),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.Issuer,
			Subject:   "umc login in",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.Timeout)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),
		},
	}
	return claims
}

func (j *JWT) CreateToken(claims *CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Key)
}

func (j *JWT) ReCreateToken(old string, claims *CustomClaims) (string, error) {
	v, err, _ := concurrency_control.Do("JWT-"+old, func() (interface{}, error) {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(j.Timeout))
		claims.BufferTime = time.Now().Add(j.BufferTime)
		claims.NotBefore = jwt.NewNumericDate(time.Now().Add(-1000))
		return j.CreateToken(claims)
	})
	return v.(string), err
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.Key, nil
	})
	if err != nil {
		return nil, err
	}
	if token != nil {
		if !token.Valid {
			return nil, TokenInvalid
		}
		if claims, ok := token.Claims.(*CustomClaims); ok {
			return claims, nil
		} else {
			return nil, TokenInvalid
		}
	}
	return nil, TokenInvalid
}
