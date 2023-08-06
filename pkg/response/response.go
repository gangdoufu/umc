package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	successStr = "success"
	errorStr   = "error"
)

type Response struct {
	Status  string
	Message string
	Data    interface{}
}

func Success(c *gin.Context) {
	c.JSON(200, &Response{Status: successStr})
}

func SuccessWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status: successStr,
		Data:   data,
	})
}

func Error(c *gin.Context, e error) {
	c.Error(e)
	c.JSON(http.StatusOK, &Response{Status: errorStr, Message: e.Error()})
}
