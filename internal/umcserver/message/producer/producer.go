package producer

import (
	"context"
	"github.com/gangdoufu/umc/internal/umcserver/sender"
)

var MessageProducer IProducer

type IProducer interface {
	SendVerificationCode(ctx context.Context, vo *sender.InfoVo) error
}
