package store

import (
	"errors"
	"github.com/gangdoufu/umc/pkg/db/model"
	"gorm.io/gorm"
)

type IBaseStore[T any] interface {
	Create(t *T) error
	CreateList(list []*T) error
	Del(id uint) error
	DelByIds(ids []uint) error
	Update(t *T) error
	UpdateList(ids []uint, t *T) error
	Save(t *T) error
	FindOne(t *T) (*T, error)
	FindById(id uint) (*T, error)
	FindList(t *T) ([]*T, error)
	FindByIds(ids []uint) ([]*T, error)
	SetDb(db *gorm.DB)
	GetDb() *gorm.DB
}

type BaseStore[T any] struct {
	DB *gorm.DB
}

func (s BaseStore[T]) SetDb(db *gorm.DB) {
	s.DB = db
}

func (s BaseStore[T]) GetDb() *gorm.DB {
	return s.DB
}

func (s BaseStore[T]) Create(t *T) error {
	return s.DB.Create(t).Error
}

func (s BaseStore[T]) CreateList(list []*T) error {
	return s.DB.CreateInBatches(list, 200).Error
}

func (s BaseStore[T]) Del(id uint) error {
	t := new(T)
	return s.DB.Delete(t, id).Error
}

func (s BaseStore[T]) DelByIds(ids []uint) error {
	t := new(T)
	return s.DB.Delete(t, ids).Error
}

func (s BaseStore[T]) Update(t *T) error {
	return s.DB.Updates(t).Error
}

func (s BaseStore[T]) UpdateList(ids []uint, t *T) error {
	return s.DB.Model(t).Where("id in ?", ids).Updates(t).Error
}

func (s BaseStore[T]) Save(t *T) error {
	return s.DB.Save(t).Error
}

func (s BaseStore[T]) FindOne(t *T) (*T, error) {
	var res T
	if err := s.DB.Where(t).First(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (s BaseStore[T]) FindById(id uint) (*T, error) {
	var res T
	if err := s.DB.Find(&res, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (s BaseStore[T]) FindList(t *T) ([]*T, error) {
	var res []*T
	if err := s.DB.Where(t).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (s BaseStore[T]) FindListByPage(t *T, page *model.PageQueryVo) (*model.PageResult[T], error) {
	var res []*T
	var count int64
	if err := s.DB.Model(new(T)).Where(t).Count(&count).Error; err != nil {
		return nil, err
	}
	if err := s.DB.Scopes(paginate(page)).Where(t).Find(&res).Error; err != nil {
		return nil, err
	}
	return model.NewPageResult(res, page, count), nil
}

func (s BaseStore[T]) FindByIds(ids []uint) ([]*T, error) {
	var res []*T
	if err := s.DB.Find(&res, ids).Error; err != nil {
		return nil, err
	}
	return res, nil
}
func paginate(page *model.PageQueryVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page.PageSize > 1000 {
			page.PageSize = 1000
		} else if page.PageSize <= 0 {
			page.PageSize = 10
		}
		offset := page.Offset()
		return db.Offset(offset).Limit(page.PageSize)
	}
}
