package vo

import (
	"github.com/gangdoufu/umc/internal/umcserver/model"
	model2 "github.com/gangdoufu/umc/pkg/db/model"
	"time"
)

type LoginType uint

const (
	AccountPassword LoginType = iota
	TelCode
	EmailCode
)

type LoginVo struct {
	Type     LoginType
	Key      string
	Value    string
	ClientId string
}

type UserGroupVo struct {
	UserId   uint
	GroupId  uint
	TenantId uint
}

type UserResourceVo struct {
	UserId       uint
	GroupId      uint
	TenantId     uint
	ResourceName string
	ResourceId   string
	ResourceCode string
}

type UserVo struct {
	UserId    uint
	CreateAt  *time.Time `json:"create_at"`
	FirstName string     `json:"first_name"` //姓
	LastName  string     `json:"last_name"`  // 名
	Nick      string     `json:"nick"`       // 昵称
	Account   string     `json:"account"`    // 账号
	Password  string     `json:"password"`   // 密码
	Birthday  string     `json:"birthday"`   // 生日
	Pic       []byte     `json:"pic"`        // 头像
	Tel       string     `json:"tel"`        // 电话
	Email     string     `json:"email"`      // 邮箱
	Status    string     `json:"status"`     // 状态 等待激活  正常 禁用 锁定
}

func (v *UserVo) GetModel() *model.User {
	user := &model.User{
		BaseModel: model2.BaseModel{ID: v.UserId, CreatedAt: v.CreateAt},
		FirstName: v.FirstName,
		LastName:  v.LastName,
		Nick:      v.Nick,
		Account:   v.Account,
		Password:  v.Password,
		Pic:       v.Pic,
		Tel:       v.Tel,
		Email:     v.Email,
		Status:    v.Status,
	}
	if v.Password != "" {
		user.SetPassword(v.Password)
	}
	if v.Birthday != "" {
		parse, err := time.Parse(time.DateOnly, v.Birthday)
		if err == nil {
			user.Birthday = &parse
		}
	}

	return user
}
