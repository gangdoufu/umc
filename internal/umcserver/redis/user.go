package redis

import (
	"context"
	"github.com/spf13/cast"
	"strings"
	"time"
)

const (
	verificationCodeKey = "verification_code"
	codeExpiration      = 5 * time.Minute // 验证码失效时间
)

const (
	codeTypeLogin = "user_login_code"
	codeTypeResetPassword
)

type UserVerificationCodeVo struct {
	UserId   uint
	Source   string
	Code     string
	CodeType string
}

func (u *UserVerificationCodeVo) GetKey() string {
	return strings.Join([]string{verificationCodeKey, cast.ToString(u.UserId), u.Source, u.CodeType}, "-")
}

func SetUserResetPasswordCode(ctx context.Context, vo *UserVerificationCodeVo) error {
	vo.Code = codeTypeResetPassword
	return SetUserVerificationCode(ctx, vo)
}

func SetUserLoginCode(ctx context.Context, vo *UserVerificationCodeVo) error {
	vo.Code = codeTypeLogin
	return SetUserVerificationCode(ctx, vo)
}

// 将验证码写入redis 默认有效时间为5分钟
func SetUserVerificationCode(ctx context.Context, vo *UserVerificationCodeVo) error {
	return client.Set(ctx, vo.GetKey(), vo.Code, codeExpiration).Err()
}

// 校验验证码 验证码正常校验之后就要删除
func CheckUserVerificationCode(ctx context.Context, vo *UserVerificationCodeVo) bool {
	res := client.Get(ctx, vo.GetKey()).String()
	if vo.Code == res {
		client.Del(ctx, vo.GetKey())
		return true
	} else {
		return false
	}
}
