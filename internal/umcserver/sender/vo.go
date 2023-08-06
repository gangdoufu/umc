package sender

import (
	"encoding/json"
	"fmt"
)

type senderType uint

type useTo uint

const (
	email senderType = iota
	tel
)
const (
	login         useTo = iota // 用于验证登录
	resetPassword              // 用于 重置密码
	activeAccount              // 用于激活账号

)

var subjectMap = map[useTo]string{
	login:         "umc 登录验证码",
	resetPassword: "umc 重置密码",
	activeAccount: "umc 账号激活",
}

type InfoVo struct {
	Type     senderType `json:"type"`
	Receiver string     `json:"receiver"`
	Info     string     `json:"info"`
	UseTo    useTo      `json:"use_to"`
}

func NewResetPasswordInfo(receiver string, infoStr string) *InfoVo {
	info := new(InfoVo)
	info.Type = email
	info.Receiver = receiver
	info.UseTo = resetPassword
	info.Info = fmt.Sprintf("请点击<a href=\"%s\">密码重置连接</a>,重新设置密码", infoStr)
	return info
}

func (v *InfoVo) String() string {
	marshal, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(marshal)
}
