package model

import "github.com/spf13/cast"

type PageQueryVo struct {
	CurPage  int
	PageSize int
}

func (v PageQueryVo) Offset() int {
	if v.CurPage <= 0 {
		v.CurPage = 1
	}
	return (v.CurPage - 1) * v.PageSize
}

type PageResult[T any] struct {
	*PageQueryVo
	Total     int64
	List      []*T
	TotalPage int64
}

func NewPageResult[T any](list []*T, v *PageQueryVo, total int64) *PageResult[T] {
	if v == nil {
		return nil
	}
	result := new(PageResult[T])
	result.Total = total
	result.TotalPage = total / cast.ToInt64(v.PageSize)
	if total%cast.ToInt64(v.PageSize) > 0 {
		result.TotalPage++
	}
	result.PageQueryVo = v
	result.List = list
	return result
}
