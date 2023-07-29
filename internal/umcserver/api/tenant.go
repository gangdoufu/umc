package api

import (
	"errors"
	"github.com/gangdoufu/umc/internal/umcserver/global"
	"github.com/gangdoufu/umc/internal/umcserver/model"
	"github.com/gangdoufu/umc/internal/umcserver/service"
	"github.com/gangdoufu/umc/internal/umcserver/service/vo"
	"github.com/gangdoufu/umc/pkg/common"
	model2 "github.com/gangdoufu/umc/pkg/db/model"
	"github.com/gangdoufu/umc/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// @Summary 创建租户
// @Produce json
// @Param app_id body string true "租户APP ID"
// @Param name body string true "租户名称"
// @Param desc body string true "说明"
// @Param name body string true "租户名称
// @Success 200 {object} model.Tenant
// @Router /manger/tenant/create [post]
func CreateTenant(c *gin.Context) {
	var tenant = &model.Tenant{}
	err := c.ShouldBind(tenant)
	if err != nil {
		if err != nil {
			response.Error(c, err)
			return
		}
	}
	db := global.DB.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	ts := service.NewTenantService(db)
	token := common.RandToken(32)
	tenant.Token = common.MD5(token)
	err = ts.CreateTenant(c.Request.Context(), tenant)
	tenant.Token = token
	if err != nil {
		response.Error(c, err)
	} else {
		response.SuccessWithData(c, tenant)
	}
}

// 更新租户信息
func UpdateTenantInfo(c *gin.Context) {
	var tenant = &model.Tenant{}
	err := c.ShouldBind(tenant)
	if err != nil {
		if err != nil {
			response.Error(c, err)
			return
		}
	}
	db := global.DB.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	ts := service.NewTenantService(db)
	err = ts.UpdateBaseInfo(c.Request.Context(), tenant)
	if err != nil {
		response.Error(c, err)
	} else {
		response.Success(c)
	}
}

// 重新生成token

func ReSetToken(c *gin.Context) {
	param := c.Param("tenant_id")
	if param == "" {
		response.Error(c, errors.New("need tenant id"))
	}
	tenantId := cast.ToUint(param)
	if tenantId == 0 {
		response.Error(c, errors.New("need tenant id"))
	}
	token := common.RandToken(32)
	tenant := &model.Tenant{
		BaseModel: model2.BaseModel{ID: tenantId},
		Token:     common.MD5(token),
	}
	db := global.DB.Begin()
	var err error
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	ts := service.NewTenantService(db)
	err = ts.UpdateBaseInfo(c.Request.Context(), tenant)
	tenant.Token = token
	if err != nil {
		response.Error(c, err)
	} else {
		response.SuccessWithData(c, tenant)
	}

}

// 新增资源

func TenantAddResource(c *gin.Context) {
	var resource = &model.Resource{}
	err := c.ShouldBind(resource)
	if err != nil {
		if err != nil {
			response.Error(c, err)
			return
		}
	}
	db := global.DB.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	ts := service.NewTenantService(db)
	err = ts.CreateResource(c.Request.Context(), resource)
	if err != nil {
		response.Error(c, err)
	} else {
		response.Success(c)
	}
}

// 新增用户组

func TenantAddGroup(c *gin.Context) {
	var group = &model.Group{}
	err := c.ShouldBind(group)
	if err != nil {
		if err != nil {
			response.Error(c, err)
			return
		}
	}
	db := global.DB.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	ts := service.NewTenantService(db)
	err = ts.CreateGroup(c.Request.Context(), group)
	if err != nil {
		response.Error(c, err)
	} else {
		response.Success(c)
	}
}

// 给用户组添加资源
func TenantGroupAddResource(c *gin.Context) {
	var groupResource = &model.GroupResource{}
	err := c.ShouldBind(groupResource)
	if err != nil {
		if err != nil {
			response.Error(c, err)
			return
		}
	}
	db := global.DB.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	ts := service.NewTenantService(db)
	err = ts.GroupAddResource(c.Request.Context(), groupResource)
	if err != nil {
		response.Error(c, err)
	} else {
		response.Success(c)
	}
}

func GetTenantBaseInfo(c *gin.Context) {
	param := c.Param("tenant_id")
	if param == "" {
		response.Error(c, errors.New("need tenant id"))
	}
	tenantId := cast.ToUint(param)
	if tenantId == 0 {
		response.Error(c, errors.New("need tenant id"))
	}
	db := global.DB
	ts := service.NewTenantService(db)
	tenant, err := ts.QueryTenantInfoById(c.Request.Context(), tenantId)
	if err != nil {
		response.Error(c, err)
	} else {
		response.SuccessWithData(c, tenant)
	}
}

func ListTenantUsers(c *gin.Context) {
	param := c.Param("tenant_id")
	if param == "" {
		response.Error(c, errors.New("need tenant id"))
	}
	tenantId := cast.ToUint(param)
	if tenantId == 0 {
		response.Error(c, errors.New("need tenant id"))
	}
	db := global.DB
	ts := service.NewTenantService(db)
	users, err := ts.ListTenantUsers(c.Request.Context(), tenantId)
	if err != nil {
		response.Error(c, err)
	} else {
		response.SuccessWithData(c, users)
	}
}

func ListTenantGroups(c *gin.Context) {
	param := c.Param("tenant_id")
	if param == "" {
		response.Error(c, errors.New("need tenant id"))
	}
	tenantId := cast.ToUint(param)
	if tenantId == 0 {
		response.Error(c, errors.New("need tenant id"))
	}
	db := global.DB
	ts := service.NewTenantService(db)
	groups, err := ts.ListTenantGroups(c.Request.Context(), tenantId)
	if err != nil {
		response.Error(c, err)
	} else {
		response.SuccessWithData(c, groups)
	}
}

func ListGroupUsers(c *gin.Context) {
	tenantParam := c.Param("tenant_id")
	groupParam := c.Param("group_id")
	if tenantParam == "" || groupParam == "" {
		response.Error(c, errors.New("need tenant id and group id"))
	}
	tenantId, groupId := cast.ToUint(tenantParam), cast.ToUint(groupParam)
	if tenantId == 0 || groupId == 0 {
		response.Error(c, errors.New("need tenant id and group id"))
	}
	db := global.DB
	ts := service.NewTenantService(db)
	users, err := ts.ListGroupUsers(c.Request.Context(), &vo.GroupVo{GroupId: groupId, TenantId: tenantId})
	if err != nil {
		response.Error(c, err)
	} else {
		response.SuccessWithData(c, users)
	}
}

func ListGroupResources(c *gin.Context) {
	groupParam := c.Param("group_id")
	if groupParam == "" {
		response.Error(c, errors.New("need  group id"))
	}
	groupId := cast.ToUint(groupParam)
	if groupId == 0 {
		response.Error(c, errors.New("need group id"))
	}
	db := global.DB
	ts := service.NewTenantService(db)
	resource, err := ts.ListGroupResource(c.Request.Context(), groupId)
	if err != nil {
		response.Error(c, err)
	} else {
		response.SuccessWithData(c, resource)
	}
}
