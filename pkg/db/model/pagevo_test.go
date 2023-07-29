package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStruct struct {
}

func TestNewPageResult(t *testing.T) {
	queryVo := PageQueryVo{
		CurPage:  1,
		PageSize: 3,
	}
	var list []*testStruct
	list = append(list, &testStruct{}, &testStruct{})
	result := NewPageResult(list, &queryVo, 29)
	var want int64 = 10
	assert.NotEmpty(t, result)
	assert.Equal(t, want, result.TotalPage)
	result = NewPageResult(list, &queryVo, 30)
	assert.Equal(t, want, result.TotalPage)
	result = NewPageResult(list, nil, 30)
	assert.Empty(t, result)
}

func TestPageQueryVo_Offset(t *testing.T) {
	query := PageQueryVo{
		CurPage:  2,
		PageSize: 3,
	}
	offset := query.Offset()
	assert.Equal(t, 3, offset)
	query.CurPage = 0
	assert.Equal(t, 0, query.Offset())
}
