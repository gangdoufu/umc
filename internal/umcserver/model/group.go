package model

import (
	model2 "github.com/gangdoufu/umc/pkg/db/model"
)

// Group 用户组
type Group struct {
	model2.BaseModel
	TenantId uint
	Name     string
	Desc     string
	Status   string
}

// GroupUser 组中包含的用户
type GroupUser struct {
	model2.BaseModel
	TenantId uint
	UserId   uint
	GroupId  uint
}

type GroupResource struct {
	model2.BaseModel
	GroupId    uint             // 用户组
	ResourceId uint             // 资源
	AuthLevel  model2.AuthLevel // 对资源的权限等级。查  增 改 删
}
