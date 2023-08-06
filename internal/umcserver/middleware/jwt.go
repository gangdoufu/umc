package middleware

import (
	"errors"
	"github.com/gangdoufu/umc/internal/umcserver/global"
	"github.com/gangdoufu/umc/internal/umcserver/redis"
	"github.com/gangdoufu/umc/internal/umcserver/service"
	mycontext "github.com/gangdoufu/umc/pkg/context"
	"github.com/gangdoufu/umc/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"time"
)

const (
	jwtToken      = "x-token"
	newJwtToken   = "x-new-token"
	tenantToken   = "x-tenant-token"
	tenantId      = "x-tenant-id"
	xUserId       = "x-user-id"
	xUserAccount  = "x-user-account"
	xUserClientId = "x-user-client-id"
)

func SetUserJwt(userId uint, clientId, account, token string, c *gin.Context) {
	c.Header(newJwtToken, token)
	c.Header(xUserId, cast.ToString(userId))
	c.Header(xUserAccount, account)
	c.Header(xUserClientId, clientId)
}

func GetToken(c *gin.Context) string {
	return c.Request.Header.Get(jwtToken)
}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		success, newToken := CheckToken(c)
		if !success {
			c.Abort()
		}
		if newToken != "" {
			c.Header(newJwtToken, newToken)
		}
		c.Next()
	}
}

func CheckToken(c *gin.Context) (success bool, newToken string) {
	success = false
	newToken = ""
	token := c.Request.Header.Get(jwtToken)
	if token == "" {
		c.JSON(http.StatusUnauthorized, "未登录或非法访问")

		return
	}
	if redis.CheckTokenInBlack(c.Request.Context(), token) {
		c.JSON(http.StatusUnauthorized, "token已失效,请重新登录")
		return
	}
	parseToken, err := global.Jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	// 如果以及到了缓冲时间.就需要重新生成一个
	if parseToken.BufferTime.Unix() <= time.Now().Unix() {
		newToken, _ = global.Jwt.ReCreateToken(token, parseToken)
		redis.AddTokenTOBlacklist(c.Request.Context(), token, parseToken.ExpiresAt.Sub(time.Now()))
	}
	success = true
	ctx := c.Request.Context()
	info := mycontext.WithUserInfo(ctx, parseToken.UserID, parseToken.Account)
	c.Request = c.Request.WithContext(info)
	return
}

func WebApiToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(tenantToken)
		tenantId := cast.ToUint(c.Request.Header.Get(tenantId))

		if token == "" || tenantId == 0 {
			response.Error(c, errors.New("请求中没有对应租户和token信息鉴权失败"))
			return
		}

		if service.CheckTenantToken(c.Request.Context(), token, tenantId) {
			response.Error(c, errors.New("token 校验失败"))
			return
		}
		c.Next()
	}
}
