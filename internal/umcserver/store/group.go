package store

import (
	"github.com/gangdoufu/umc/internal/umcserver/model"
	"github.com/gangdoufu/umc/pkg/db/store"
	"gorm.io/gorm"
)

type IGroupStore interface {
	store.IBaseStore[model.Group]
	QueryTenantGroupList(tenantId uint) ([]*model.Group, error)
	QueryTenantGroupNyName(tenantId uint, groupName string) (*model.Group, error)
	QueryTenantDefaultGroup(tenantId uint) (*model.Group, error)
	CreateGuestGroup(tenantId uint) (*model.Group, error)
	QueryUserGroupList(tenantId, userId uint) ([]*model.Group, error)
	CreateMangerGroup(tenantId uint) (*model.Group, error)
	CreateDefaultGroups(tenantId uint) ([]*model.Group, error)
}

type groupStore struct {
	store.BaseStore[model.Group]
}

func (g groupStore) CreateDefaultGroups(tenantId uint) ([]*model.Group, error) {
	var groups = []*model.Group{&model.Group{
		TenantId: tenantId,
		Name:     "manger",
	}, &model.Group{
		TenantId: tenantId,
		Name:     "guest",
	}}
	if err := g.GetDb().Create(groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (g groupStore) CreateMangerGroup(tenantId uint) (*model.Group, error) {
	group := &model.Group{
		TenantId: tenantId,
		Name:     "manger",
	}
	err := g.Create(group)
	if err != nil {
		return nil, err
	}
	return group, nil
}
func (g groupStore) CreateGuestGroup(tenantId uint) (*model.Group, error) {
	group := &model.Group{
		TenantId: tenantId,
		Name:     "guest",
	}
	err := g.Create(group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (g groupStore) QueryTenantDefaultGroup(tenantId uint) (*model.Group, error) {
	return g.QueryTenantGroupNyName(tenantId, "guest")
}

func (g groupStore) QueryTenantGroupNyName(tenantId uint, groupName string) (*model.Group, error) {
	var res model.Group
	if err := g.GetDb().Model(&model.Group{}).Where(&model.Group{TenantId: tenantId, Name: groupName}).Find(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

func (g groupStore) QueryTenantGroupList(tenantId uint) ([]*model.Group, error) {
	var res []*model.Group
	if err := g.GetDb().Model(&model.Group{}).Where(&model.Group{TenantId: tenantId}).Find(res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// 查询用户的所有群组
func (g groupStore) QueryUserGroupList(tenantId, userId uint) ([]*model.Group, error) {
	var res []*model.Group
	db := g.GetDb()
	if err := db.Model(&model.Group{}).Where("id in (?)",
		db.Model(&model.GroupUser{}).Select("group_id").
			Where("user_id = ? and tenant_id=?", userId, tenantId)).Find(res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func NewGroupStore(db *gorm.DB) IGroupStore {
	group := groupStore{}
	group.DB = db
	return group
}

type IGroupUsersStore interface {
	store.IBaseStore[model.GroupUser]
	ListTenantUserIds(tenantId uint) ([]uint, error)
	RemoveUserGroup(userId, groupId uint) error
	RemoveUserTenant(userId, tenantId uint) error
	ListGroupUserIds(tenantId, groupId uint) ([]uint, error)
}

type groupUsersStore struct {
	store.BaseStore[model.GroupUser]
}

func (gu groupUsersStore) RemoveUserTenant(userId, tenantId uint) error {
	return gu.GetDb().Delete(&model.GroupUser{UserId: userId, TenantId: tenantId}).Error
}
func (gu groupUsersStore) RemoveUserGroup(userId, groupId uint) error {
	return gu.GetDb().Delete(&model.GroupUser{UserId: userId, GroupId: groupId}).Error
}

func (gu groupUsersStore) UserAddGroups(tenantId, userId uint, groupIds []uint) error {
	groupUsers := make([]*model.GroupUser, len(groupIds))
	for _, id := range groupIds {
		groupUsers = append(groupUsers, &model.GroupUser{
			TenantId: tenantId,
			UserId:   userId,
			GroupId:  id,
		})
	}
	return gu.GetDb().Create(groupUsers).Error
}

func NewGroupUsersStore(db *gorm.DB) IGroupUsersStore {
	group := groupUsersStore{}
	group.DB = db
	return group
}
func (gu groupUsersStore) ListTenantUserIds(tenantId uint) ([]uint, error) {
	var users []uint
	if err := gu.GetDb().Model(&model.GroupUser{}).
		Where(&model.GroupUser{TenantId: tenantId}).Select("user_id").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil

}

func (gu groupUsersStore) ListGroupUserIds(tenantId, groupId uint) ([]uint, error) {
	var users []uint
	if err := gu.GetDb().Model(&model.GroupUser{}).
		Where(&model.GroupUser{TenantId: tenantId, GroupId: groupId}).Select("user_id").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

type IGroupResourceStore interface {
	store.IBaseStore[model.GroupResource]
}

type groupResourceStore struct {
	store.BaseStore[model.GroupResource]
}

func NewGroupResourceStore(db *gorm.DB) IGroupResourceStore {
	group := groupResourceStore{}
	group.DB = db
	return group
}
