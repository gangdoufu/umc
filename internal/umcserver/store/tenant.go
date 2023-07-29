package store

import (
	"github.com/gangdoufu/umc/internal/umcserver/model"
	"github.com/gangdoufu/umc/pkg/db/store"
	"gorm.io/gorm"
)

type ITenantStore interface {
	store.IBaseStore[model.Tenant]
	ListUserTenants(userId uint) ([]*model.TenantShowVo, error)
}

type tenantStore struct {
	store.BaseStore[model.Tenant]
}

// 查询用户的所有
func (s tenantStore) ListUserTenants(userId uint) ([]*model.TenantShowVo, error) {
	var list []*model.TenantShowVo
	db := s.GetDb()
	if err := db.Model(&model.Tenant{}).Where("id in (?)",
		db.Model(&model.Group{}).Select("tenant_id").Where("id in (?)",
			db.Model(&model.GroupUser{}).Select("group_id").
				Where("user_id = ?", userId))).Find(list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func NewTenantStore(db *gorm.DB) ITenantStore {
	tenant := &tenantStore{}
	tenant.DB = db
	return tenant
}
