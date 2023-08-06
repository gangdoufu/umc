package producer

import (
	"context"
	"github.com/gangdoufu/umc/internal/umcserver/message"
	"github.com/gangdoufu/umc/internal/umcserver/sender"
	"github.com/go-redis/redis/v8"
)

type RedisQueue struct {
	Client redis.UniversalClient
}

func (r *RedisQueue) SendVerificationCode(ctx context.Context, vo *sender.InfoVo) error {
	return r.Client.Do(ctx, "publish", message.VerificationCodeKey, vo.String()).Err()
}

func NewREdisQueue(client redis.UniversalClient) {

}
