package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const xRequestId = "X-Request-ID"

func RequestId(c *gin.Context) {
	rid := c.GetHeader(xRequestId)
	if rid == "" {
		rid = uuid.New().String()
		c.Request.Header.Add(xRequestId, rid)
	}
	c.Header(xRequestId, rid)
	c.Next()
}
