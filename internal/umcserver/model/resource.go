package model

import (
	"github.com/gangdoufu/umc/pkg/db/model"
)

type Resource struct {
	model.BaseModel
	TenantId uint   `json:"tenant_id"` // 所属租户
	Name     string `json:"name"`      // 名称
	Code     string `json:"code"`      // code
	Desc     string `json:"desc"`      // 说明
	Type     string `json:"type"`      // 资源类型
}
