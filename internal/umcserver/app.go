package umcserver

import (
	"github.com/gangdoufu/umc/internal/umcserver/config"
	"github.com/gangdoufu/umc/internal/umcserver/global"
	"github.com/gangdoufu/umc/internal/umcserver/initialize"
	"github.com/gangdoufu/umc/internal/umcserver/redis"
	"github.com/gangdoufu/umc/internal/umcserver/router"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"strconv"
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

	if global.Config.Server.UseSSL {
		ssl := global.Config.SSL
		a.Engin.Use(TlsHandler(8088))
		a.Engin.RunTLS(":"+strconv.Itoa(8088), ssl.CertFile, ssl.PrivateKeyFile)
	} else {
		a.Engin.Run("0.0.0.0:8088")
	}

}

func TlsHandler(port int) gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ":" + strconv.Itoa(port),
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
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
