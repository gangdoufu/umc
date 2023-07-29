package initialize

import (
	"github.com/gangdoufu/umc/internal/umcserver/config"
	"github.com/gangdoufu/umc/pkg/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitLogger(logConfig *config.Log) *zap.Logger {
	op := &log.Option{
		Path:         logConfig.Path,
		MaxAge:       logConfig.MaxAge,
		MaxSize:      logConfig.MaxSize,
		Level:        logConfig.Level,
		LogInConsole: logConfig.LogIncConsole,
		Format:       logConfig.Format,
		LevelEncode:  logConfig.LevelEncode,
	}
	return log.Build(op)
}

func SetGormLogger(db *gorm.DB) {

	gormLogger := log.NewGormLogger(log.Logger())
	db.Logger = gormLogger
}
