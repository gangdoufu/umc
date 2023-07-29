package router

import (
	"github.com/gangdoufu/umc/internal/umcserver/api"
	"github.com/gin-gonic/gin"
)

func initUserRouter(router *gin.RouterGroup) {
	user := router.Group("/user")
	{
		user.POST("create", api.CreateUser)
		user.POST("add_tenant", api.UserAddInTenant)
		user.POST("add_group", api.UserAddInGroup)
		user.GET("/info/:userid", api.GetUserBaseInfo)
		user.GET("/info/tenants/:userid", api.GetUserTenants)
		user.POST("/info/groups", api.GetUserGroups)
	}
}

func InitUserWebApiRouter(router *gin.RouterGroup) {

}

// 用户

// 登录

// 注册

// 管理员邀请用户

// 用户注册

// 用户邮箱激活

// 用户

// 管理员禁用用户

// 用户更新信息
