package middleware

import (
	"github.com/gangdoufu/umc/internal/umcserver/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	ginContextTransaction = "umc-gin-transaction"
)

func Transaction(logger *zap.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Next()
		db := getGinDB(c)
		if db != nil {
			if len(c.Errors) > 0 {
				GetLogger(c, logger).Info("request error db rollback", zap.Error(c.Err()))
				db.Rollback()
			} else {
				db.Commit()
			}
			c.Set(ginContextTransaction, nil)
		}
	}
}

func getGinDB(c *gin.Context) *gorm.DB {
	value, exists := c.Get(ginContextTransaction)
	if exists {
		if v, ok := value.(*gorm.DB); ok {
			return v
		}
	}
	return nil
}

func SetGinDB(transaction bool, c *gin.Context) *gorm.DB {
	if transaction {
		db := global.DB.Begin()
		c.Set(ginContextTransaction, db)
		return db
	}
	return global.DB
}
