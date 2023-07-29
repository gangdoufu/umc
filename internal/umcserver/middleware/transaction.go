package middleware

import (
	"github.com/gangdoufu/umc/pkg/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewTransaction() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Next()
		value, exists := c.Get(common.GinUseGlobalTransaction)
		if exists {
			if db, ok := value.(*gorm.DB); ok {
				if c.Err() != nil {
					db.Rollback()
				} else {
					db.Commit()
				}
			}

		}

	}
}
func SetTransaction(c *gin.Context, db *gorm.DB) {
	c.Set(common.GinUseGlobalTransaction, db)
}

//func GetGinDB(global bool) *gorm.DB {
//
//}
