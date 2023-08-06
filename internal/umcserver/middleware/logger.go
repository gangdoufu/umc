package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"io"
	"time"
)

func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path, method := c.Request.URL.Path, c.Request.Method
		clientId := c.GetHeader("X-Forwarded-For")
		cLogger := GetLogger(c, logger)
		var inPutValues []zap.Field
		params := c.Params
		if len(params) > 0 {
			for _, param := range params {
				inPutValues = append(inPutValues, zap.String(param.Key, param.Value))
			}
		}
		var bodyInfo []byte
		bodyInfo, err := c.GetRawData()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyInfo))
		if err == nil {
			inPutValues = append(inPutValues, zap.String("request body", string(bodyInfo)))
		}
		cLogger.Info("request start,",
			zap.String("client-id", clientId),
			zap.String("path", path),
			zap.String("method", method))
		if len(inPutValues) > 0 {

			cLogger.Info("request input param", inPutValues...)
		}
		c.Next()
		elapsed := time.Since(start)
		cLogger.Info("request end,", zap.Duration("elapsed", elapsed))
	}
}

func GetLogger(c *gin.Context, logger *zap.Logger) *zap.Logger {
	rid := c.GetHeader(xRequestId)
	value, exists := c.Get("user_account")
	var fields []zap.Field
	fields = append(fields, zap.String("request_id", rid))
	if exists {
		fields = append(fields, zap.String("account", cast.ToString(value)))
	}
	cLogger := logger.With(fields...)
	return cLogger
}
