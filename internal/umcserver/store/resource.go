package store

import (
	"github.com/gangdoufu/umc/internal/umcserver/model"
	"github.com/gangdoufu/umc/pkg/db/store"
	"gorm.io/gorm"
)

type IResourceStore interface {
	store.IBaseStore[model.Resource]
}

type resourceStore struct {
	store.BaseStore[model.Resource]
}

func NewResourceStore(db *gorm.DB) IResourceStore {
	r := &resourceStore{}
	r.DB = db
	return r
}
