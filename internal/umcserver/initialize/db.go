package initialize

import (
	"fmt"
	"github.com/gangdoufu/umc/internal/umcserver/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func InitDb(mysqlConfig *config.Mysql) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		mysqlConfig.Username, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	s, err := db.DB()
	if err != nil {
		panic(err)
	}
	s.SetMaxOpenConns(mysqlConfig.MaxOpenConnections)
	s.SetMaxIdleConns(mysqlConfig.MaxIdleConnections)
	s.SetConnMaxLifetime(time.Duration(mysqlConfig.MaxConnectionLifeTime) * time.Second)
	return db
}
