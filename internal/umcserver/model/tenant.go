package model

import (
	"github.com/gangdoufu/umc/pkg/db/model"
)

// Tenant 租户
type Tenant struct {
	model.BaseModel
	ParentId uint   // 父租户ID
	AppId    string // 应用Id
	Host     string // 租户后续访问的主机.用于跨域配置
	Name     string // 名称
	Desc     string // 说明
	Manger   uint   // 管理员
	Token    string // token 在创建的时候生成,后台存的是md5.需要在创建的时候用户保存.如果丢失了只能由管理员重新生成
}

type TenantShowVo struct {
	ID    uint
	AppId string
	Name  string
}
