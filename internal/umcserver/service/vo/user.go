package vo

type LoginType uint

const (
	AccountPassword LoginType = iota
	TelCode
	EmailCode
)

type LoginVo struct {
	Type  LoginType
	Key   string
	Value string
}

type UserGroupVo struct {
	UserId   uint
	GroupId  uint
	TenantId uint
}
