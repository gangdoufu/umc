package umcserver

import (
	"github.com/gangdoufu/umc/internal/umcserver/config"
	"github.com/gangdoufu/umc/internal/umcserver/global"
	"github.com/gangdoufu/umc/internal/umcserver/initialize"
	"github.com/gangdoufu/umc/internal/umcserver/redis"
	"github.com/gangdoufu/umc/internal/umcserver/router"
	"github.com/gin-gonic/gin"
)

type UMCApp struct {
	Engin *gin.Engine
}

func NewUMCApp() *UMCApp {
	return &UMCApp{Engin: gin.Default()}
}

func (a *UMCApp) Run() {
	global.Config = config.LoadConfig()
	InitCommon(global.Config)
	router.InitRouter(a.Engin)
	a.Engin.Run("0.0.0.0:8088")
}

func InitCommon(config *config.Config) {
	// log
	global.Logger = initialize.InitLogger(&config.Log)
	// db
	global.DB = initialize.InitDb(&config.Mysql)
	initialize.SetGormLogger(global.DB)
	// redis
	redis.NewRedisCache(&redis.Option{Addr: config.Redis.Host, Password: config.Redis.Password})
	// queue
	// jwt
	global.Jwt = initialize.InitJwt(&config.Jwt)
}
