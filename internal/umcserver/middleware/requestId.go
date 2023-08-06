package middleware

import (
	mycontext "github.com/gangdoufu/umc/pkg/context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

const xRequestId = "X-Request-ID"

func RequestId(c *gin.Context) {
	rid := c.GetHeader(xRequestId)
	if rid == "" {
		rid = uuid.New().String()
		c.Request.Header.Add(xRequestId, rid)
	}
	c.Header(xRequestId, rid)
	ctx := c.Request.Context()
	ctx = mycontext.WithRequestInfo(ctx, rid, time.Now())
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}
