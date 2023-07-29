package context

import (
	"gorm.io/gorm"
	"time"
)

type UserVo struct {
	UserId  uint
	Account string
}

type RequestVo struct {
	RequestId string
	RequestAt time.Time
}

type TransactionVo struct {
	Db *gorm.DB
}
