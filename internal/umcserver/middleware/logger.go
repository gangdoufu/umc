package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path, method := c.Request.URL.Path, c.Request.Method
		clientId := c.GetHeader("X-Forwarded-For")
		c.Next()
		elapsed := time.Since(start)
		var fields []zap.Field
		fields = append(fields, zap.Duration("elapsed", elapsed),
			zap.String("path", path),
			zap.String("method", method),
			zap.String("client-id", clientId),
		)
	}
}
