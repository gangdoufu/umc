package router

import (
	"github.com/gangdoufu/umc/internal/umcserver/api"
	"github.com/gin-gonic/gin"
)

func initUserRouter(router *gin.RouterGroup) {
	user := router.Group("/user")
	{
		user.POST("/create", api.CreateUser)
		user.POST("/add_tenant", api.UserAddInTenant)
		user.POST("/add_group", api.UserAddInGroup)
		// 将用户从租户中移除
		user.POST("/remove_tenant", api.UserRemoveTenant)
		// 将用户从用户组移除
		user.POST("/remove_group", api.UserRemoveGroup)

		user.GET("/info/:userid", api.GetUserBaseInfo)
		user.GET("/info/tenants/:userid", api.GetUserTenants)
		user.POST("/info/groups", api.GetUserGroups)
		// 忘记密码
		user.GET("/forget_password")
		user.POST("/reset_password")

		// 修改密码
		user.POST("/change_password")

	}

}

func InitUserWebApiRouter(router *gin.RouterGroup) {
	user := router.Group("/user")
	{
		// 校验用户权限
		user.POST("/check_auth")
		//
		user.POST("/check_jwt")
		// 获取 用户的用户组
		user.GET("/list_group/:user_id/:tenant_id")
		// 获取用户组的资源
		user.GET("/list_resources/:user_id/:group_id")

	}
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
