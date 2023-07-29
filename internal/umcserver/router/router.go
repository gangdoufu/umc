package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func InitRouter(router *gin.Engine) {
	manger := router.Group("/manger")
	{
		initUserRouter(manger)
		initTenantManger(manger)
		manger.GET("/heartbeats", func(c *gin.Context) {
			c.JSON(http.StatusOK, time.Now())
		})
	}
}

func InitWebApi(router *gin.Engine) {
	webApi := router.Group("/api")
	{
		InitUserWebApiRouter(webApi)
		initTenantWebApiRouter(webApi)
	}
}
