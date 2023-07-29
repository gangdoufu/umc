package vo

type TenantVo struct {
	TenantId uint
	UserIds  []uint
	GroupIds []uint
}
type GroupVo struct {
	TenantId uint
	GroupId  uint
	UserIds  []uint
}
