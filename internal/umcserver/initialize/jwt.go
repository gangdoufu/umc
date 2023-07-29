package initialize

import (
	"github.com/gangdoufu/umc/internal/umcserver/config"
	"github.com/gangdoufu/umc/pkg/common"
)

func InitJwt(jwt *config.Jwt) *common.JWT {
	return common.NewJWT(jwt.Key, jwt.Timeout, jwt.Issuer, jwt.BufferTime)
}
