package model

import (
	"github.com/gangdoufu/umc/pkg/db/model"
)

// Tenant 租户
type Tenant struct {
	model.BaseModel
	ParentId uint   `json:"parent_id"` // 父租户ID
	AppId    string `json:"app_id"`    // 应用Id
	Host     string `json:"host"`      // 租户后续访问的主机.用于跨域配置
	Name     string `json:"name"`      // 名称
	Desc     string `json:"desc"`      // 说明
	Manger   uint   `json:"manger"`    // 管理员
	Token    string `json:"token"`     // token 在创建的时候生成,后台存的是md5.需要在创建的时候用户保存.如果丢失了只能由管理员重新生成
}

type TenantShowVo struct {
	ID    uint   `json:"id"`
	AppId string `json:"app_id"`
	Name  string `json:"name"`
}
