package store

import (
	"errors"
	"github.com/gangdoufu/umc/internal/umcserver/model"
	"github.com/gangdoufu/umc/pkg/db/store"
	"gorm.io/gorm"
)

type IUserStore interface {
	store.IBaseStore[model.User]
	QueryByAccount(account string) (*model.User, error)
	QueryUserPasswordByAccount(account string) (*model.UserPassword, error)
}

type userStore struct {
	store.BaseStore[model.User]
}

func (u userStore) QueryUserPasswordByAccount(account string) (*model.UserPassword, error) {
	var up = &model.UserPassword{}
	if err := u.GetDb().Model(&model.User{}).Where(&model.User{Account: account}).First(up).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		return nil, err
	}
	return up, nil
}

// 依据账号查询信息
func (u userStore) QueryByAccount(account string) (*model.User, error) {
	var res *model.User
	if err := u.GetDb().Model(&model.User{}).Where(&model.User{Account: account}).First(res).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return res, nil
}

func NewUserStore(db *gorm.DB) IUserStore {
	user := &userStore{}
	user.DB = db
	return user
}

type IUserClientStore interface {
	store.IBaseStore[model.UserClient]
}

type userClientStore struct {
	store.BaseStore[model.UserClient]
}

func NewUserClientStore(db *gorm.DB) IUserClientStore {
	user := &userClientStore{}
	user.DB = db
	return user
}

type IUserPasswordHistoryStore interface {
	store.IBaseStore[model.UserPasswordHistory]
}

type userPasswordHistoryStore struct {
	store.BaseStore[model.UserPasswordHistory]
}

func NewUserPasswordHistoryStore(db *gorm.DB) IUserPasswordHistoryStore {
	user := &userPasswordHistoryStore{}
	user.DB = db
	return user
}
