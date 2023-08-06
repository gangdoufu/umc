package model

import (
	"github.com/gangdoufu/umc/pkg/common"
	"github.com/gangdoufu/umc/pkg/db/model"
	"gorm.io/gorm"
	"time"
)

type User struct {
	model.BaseModel
	FirstName string     `json:"first_name"` //姓
	LastName  string     `json:"last_name"`  // 名
	Nick      string     `json:"nick"`       // 昵称
	Account   string     `json:"account"`    // 账号
	Password  string     `json:"password"`   // 密码
	Birthday  *time.Time `json:"birthday"`   // 生日
	Pic       []byte     `json:"pic"`        // 头像
	Tel       string     `json:"tel"`        // 电话
	Email     string     `json:"email"`      // 邮箱
	Status    string     `json:"status"`     // 状态 等待激活  正常 禁用 锁定
}

// 查询的返回不带密码
func (u *User) AfterFind(tx *gorm.DB) error {
	u.Password = ""
	return nil
}

type UserShowVo struct {
	model.BaseModel
	FirstName string     `json:"first_name"` //姓
	LastName  string     `json:"last_name"`  // 名
	Nick      string     `json:"nick"`       // 昵称
	Account   string     `json:"account"`    // 账号
	Birthday  *time.Time `json:"birthday"`   // 生日
	Pic       []byte     `json:"pic"`        // 头像
	Tel       string     `json:"tel"`        // 电话
	Email     string     `json:"email"`      // 邮箱
	Status    string     `json:"status"`     // 状态 等待激活  正常 禁用 锁定
}

type UserLoginVo struct {
	ID       uint   `json:"id"`
	Account  string `json:"account"`
	ClientId string `json:"client_id"`
	Token    string `json:"token"`
}
type UserPassword struct {
	ID       uint   `json:"id"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

func (u User) GetStatus() string {
	return u.Status
}

func (u User) SetUserStatusNormal() {
	u.Status = "normal"
}

func (u User) SetUserStatusWaitingActive() {
	u.Status = "wait_active"
}

func (u User) SetUserStatusLocking() {
	u.Status = "locking"
}

func (u User) SetUserStatusDisable() {
	u.Status = "disable"
}

func (u User) UserIsWaitingActive() bool {
	return u.Status == "wait_active"
}

func (u User) UserIsLocking() bool {
	return u.Status == "locking"
}

func (u User) UserIsDisable() bool {
	return u.Status == "disable"
}

func (u User) UserIsNormal() bool {
	return u.Status == "normal"
}

// UserClient 用户的记住d
type UserClient struct {
	model.BaseModel
	ClientId   string     `json:"client_id"` // ID client 的唯一编码在登录
	ClientName string     `json:"client_name"`
	Status     uint8      `json:"status"`
	HeartBeat  *time.Time `json:"heart_beat"`
	ClientType string     `json:"client_type"`
}

// UserPasswordHistory 用户用过的密码 用于用户密码强度校验.在一定时间或者次数内不可重复
type UserPasswordHistory struct {
	model.BaseModel
	UserId   uint   `json:"user_id"`
	Password string `json:"password"`
}

func (u *User) SetPassword(str string) {
	u.Password = common.MD5(str)
}

func (u *User) CheckPassword(str string) bool {
	return u.Password == common.MD5(str)
}
