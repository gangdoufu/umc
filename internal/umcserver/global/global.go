package global

import (
	"github.com/gangdoufu/umc/internal/umcserver/config"
	"github.com/gangdoufu/umc/pkg/common"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Config *config.Config
var Logger *zap.Logger
var Jwt *common.JWT
