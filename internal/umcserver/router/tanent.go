package router

import (
	"github.com/gangdoufu/umc/internal/umcserver/api"
	"github.com/gin-gonic/gin"
)

func initTenantManger(router *gin.RouterGroup) {
	tenant := router.Group("/tenant")
	{
		// 创建租户
		tenant.POST("/create", api.CreateTenant)
		// 更新租户信息
		tenant.POST("/update", api.UpdateTenantInfo)
		// 重置token
		tenant.GET("/reset_token", api.ReSetToken)
		// 租户增加用户组
		tenant.POST("/group/create", api.TenantAddGroup)
		// 用户组增加资源权限
		tenant.POST("/group/add_resource", api.TenantGroupAddResource)
		// 租户添加资源
		tenant.POST("/resource/create", api.TenantAddResource)
		// 查询租户基本信息
		tenant.GET("/get/:tenant_id", api.GetTenantBaseInfo)
		// 获取租户下面的所有用户
		tenant.GET("/list_users/:tenant_id", api.ListTenantUsers)
		// 获取租户的所有用户组
		tenant.GET("/list_group/:tenant_id", api.ListTenantGroups)
		// 获取用户组下的所有用户
		tenant.GET("/list_group_user/:tenant_id/:group_id", api.ListGroupUsers)
		// 获取用户组下所有资源
		tenant.GET("/list_group_resource/:tenant_id/:group_id", api.ListGroupResources)
	}
}

func initTenantWebApiRouter(router *gin.RouterGroup) {

}
