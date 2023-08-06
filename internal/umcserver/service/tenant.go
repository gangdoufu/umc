package service

import (
	"context"
	"github.com/gangdoufu/umc/internal/umcserver/global"
	"github.com/gangdoufu/umc/internal/umcserver/model"
	"github.com/gangdoufu/umc/internal/umcserver/redis"
	"github.com/gangdoufu/umc/internal/umcserver/service/vo"
	"github.com/gangdoufu/umc/internal/umcserver/store"
	"github.com/gangdoufu/umc/pkg/common"
	"gorm.io/gorm"
)

type ITenantService interface {
	CreateTenant(ctx context.Context, tenant *model.Tenant) error
	UpdateBaseInfo(ctx context.Context, tenant *model.Tenant) error
	UpdateToken(ctx context.Context, tenant *model.Tenant) error
	CreateGroup(ctx context.Context, group *model.Group) error
	CreateResource(ctx context.Context, resource *model.Resource) error
	GroupAddResource(ctx context.Context, groupResource *model.GroupResource) error
	QueryTenantInfoById(ctx context.Context, tenantId uint) (*model.Tenant, error)
	ListTenantUsers(ctx context.Context, tenantId uint) (*vo.TenantVo, error)
	ListGroupUsers(ctx context.Context, groupVo *vo.GroupVo) (*vo.GroupVo, error)
	ListTenantGroups(ctx context.Context, tenantId uint) ([]*model.Group, error)
	ListGroupResource(ctx context.Context, groupId uint) ([]*model.GroupResource, error)
}
type tenantService struct {
	factory store.IFactory
}

func NewTenantService(db *gorm.DB) ITenantService {
	return tenantService{factory: store.NewFactory(db)}
}

// 创建租户.创建租户 每个租户创建都需要创建一个默认的guest组 和一个manger 组。租户中所有的人员都会加入到guest组.创建的时候需要指定管理员,
// 如果没有指定则创建人是默认管理员,管理员后续也可以修改.但是需要至少一个人

func (s tenantService) CreateTenant(ctx context.Context, tenant *model.Tenant) error {
	tenants := s.factory.Tenants()
	groups := s.factory.Groups()
	//users := s.factory.GroupUsers()
	err := tenants.Create(tenant)
	if err != nil {
		return err
	}
	_, err = groups.CreateDefaultGroups(tenant.ID)

	if err != nil {
		return err
	}

	return nil
}
func (s tenantService) UpdateBaseInfo(ctx context.Context, tenant *model.Tenant) error {
	tenants := s.factory.Tenants()
	updateTenant := &model.Tenant{Desc: tenant.Desc}
	updateTenant.ID = tenant.ID
	return tenants.Update(updateTenant)
}

func (s tenantService) UpdateToken(ctx context.Context, tenant *model.Tenant) error {
	tenants := s.factory.Tenants()
	updateTenant := &model.Tenant{Token: tenant.Token}
	updateTenant.ID = tenant.ID
	return tenants.Update(updateTenant)
}

// 查询租户详细信息
func (s tenantService) QueryTenantInfoById(ctx context.Context, tenantId uint) (*model.Tenant, error) {
	tenants := s.factory.Tenants()
	return tenants.FindById(tenantId)
}

func (s tenantService) CreateGroup(ctx context.Context, group *model.Group) error {
	groups := s.factory.Groups()
	return groups.Create(group)
}

func (s tenantService) CreateResource(ctx context.Context, resource *model.Resource) error {
	resources := s.factory.Resources()
	return resources.Create(resource)
}
func (s tenantService) GroupAddResource(ctx context.Context, groupResource *model.GroupResource) error {
	resources := s.factory.GroupResources()
	return resources.Create(groupResource)
}
func (s tenantService) ListTenantGroups(ctx context.Context, tenantId uint) ([]*model.Group, error) {
	groups := s.factory.Groups()
	groupList, err := groups.FindList(&model.Group{TenantId: tenantId})
	if err != nil {
		return nil, err
	}
	return groupList, err
}

func (s tenantService) ListTenantUsers(ctx context.Context, tenantId uint) (*vo.TenantVo, error) {
	users := s.factory.GroupUsers()
	ids, err := users.ListTenantUserIds(tenantId)
	if err != nil {
		return nil, err
	}
	return &vo.TenantVo{TenantId: tenantId, UserIds: ids}, nil
}

func (s tenantService) ListGroupUsers(ctx context.Context, groupVo *vo.GroupVo) (*vo.GroupVo, error) {
	users := s.factory.GroupUsers()
	ids, err := users.ListGroupUserIds(groupVo.TenantId, groupVo.GroupId)
	if err != nil {
		return nil, err
	}
	groupVo.UserIds = ids
	return groupVo, nil
}

func (s tenantService) ListGroupResource(ctx context.Context, groupId uint) ([]*model.GroupResource, error) {
	resources := s.factory.GroupResources()
	return resources.FindList(&model.GroupResource{GroupId: groupId})
}

func CheckTenantToken(ctx context.Context, token string, tenantId uint) bool {

	tenantStore := store.NewTenantStore(global.DB)
	tenant, err := tenantStore.FindById(tenantId)
	if err != nil {
		return false
	}
	if tenant.Token != common.MD5(token) {
		return redis.CheckTenantDynamicToken(ctx, tenant.AppId, tenant.Token)
	}
	return true
}
