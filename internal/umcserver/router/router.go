package router

import (
	"github.com/gangdoufu/umc/internal/umcserver/api"
	"github.com/gangdoufu/umc/internal/umcserver/global"
	"github.com/gangdoufu/umc/internal/umcserver/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func InitRouter(router *gin.Engine) {
	router.Use(middleware.RequestId)
	router.Use(middleware.Transaction(global.Logger))
	router.Use(middleware.Logger(global.Logger))
	// 登录
	router.POST("/login", api.Login)
	// 注销
	router.POST("/logout", api.LoginOut)
	// 注册
	router.POST("/register", api.Register)

	manger := router.Group("/manger")
	{
		manger.Use(middleware.JwtAuth())
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
		webApi.Use(middleware.WebApiToken())
		InitUserWebApiRouter(webApi)
		initTenantWebApiRouter(webApi)
	}
}
