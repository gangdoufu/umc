package model

import (
	"github.com/gangdoufu/umc/pkg/db/model"
)

type Resource struct {
	model.BaseModel
	TenantId uint   // 所属租户
	Name     string // 名称
	Desc     string // 说明
	Type     string // 资源类型
}
