package middleware

import (
	"github.com/gangdoufu/umc/internal/umcserver/global"
	"github.com/gangdoufu/umc/internal/umcserver/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	jwtToken    = "x-token"
	newJwtToken = "x-new-token"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(jwtToken)
		if token == "" {
			c.JSON(http.StatusUnauthorized, "未登录或非法访问")
			c.Abort()
			return
		}
		if redis.CheckTokenInBlack(c.Request.Context(), token) {
			c.JSON(http.StatusUnauthorized, "token已失效,请重新登录")
			c.Abort()
			return
		}
		parseToken, err := global.Jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err)
			c.Abort()
			return
		}
		// 如果以及到了缓冲时间.就需要重新生成一个
		if parseToken.BufferTime.Unix() <= time.Now().Unix() {
			newToken, _ := global.Jwt.ReCreateToken(token, parseToken)
			c.Header(newJwtToken, newToken)
			redis.AddTokenTOBlacklist(c.Request.Context(), token, parseToken.ExpiresAt.Sub(time.Now()))
		}
		c.Next()
	}
}
