package model

import (
	model2 "github.com/gangdoufu/umc/pkg/db/model"
)

// Group 用户组
type Group struct {
	model2.BaseModel
	TenantId uint   `json:"tenant_id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Status   string `json:"status"`
}

// GroupUser 组中包含的用户
type GroupUser struct {
	model2.BaseModel
	TenantId uint `json:"tenant_id"`
	UserId   uint `json:"user_id"`
	GroupId  uint `json:"group_id"`
}

type GroupResource struct {
	model2.BaseModel
	GroupId    uint `json:"group_id"`    // 用户组
	ResourceId uint `json:"resource_id"` // 资源
	AuthLevel  int  `json:"auth_level"`  // 对资源的权限等级。查  增 改 删
}
