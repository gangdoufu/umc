package store

import "gorm.io/gorm"

type IFactory interface {
	Users() IUserStore
	UserClients() IUserClientStore
	UserPasswordHistory() IUserPasswordHistoryStore
	Groups() IGroupStore
	GroupUsers() IGroupUsersStore
	GroupResources() IGroupResourceStore
	Tenants() ITenantStore
	Resources() IResourceStore
}

type factory struct {
	db *gorm.DB
}

func NewFactory(db *gorm.DB) IFactory {
	return &factory{db: db}
}

func (f factory) Users() IUserStore {
	return NewUserStore(f.db)
}

func (f factory) UserClients() IUserClientStore {
	return NewUserClientStore(f.db)
}

func (f factory) UserPasswordHistory() IUserPasswordHistoryStore {
	return NewUserPasswordHistoryStore(f.db)
}

func (f factory) Groups() IGroupStore {
	return NewGroupStore(f.db)
}

func (f factory) GroupUsers() IGroupUsersStore {
	return NewGroupUsersStore(f.db)
}

func (f factory) GroupResources() IGroupResourceStore {
	return NewGroupResourceStore(f.db)
}

func (f factory) Tenants() ITenantStore {
	return NewTenantStore(f.db)
}

func (f factory) Resources() IResourceStore {
	return NewResourceStore(f.db)
}
