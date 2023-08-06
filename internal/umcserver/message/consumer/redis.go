package consumer

import (
	"context"
	"github.com/gangdoufu/umc/internal/umcserver/global"
	"github.com/go-redis/redis/v8"
)

var logger = global.Logger.Named("redis-qune")

type redisConsumer struct {
	Client redis.UniversalClient
}

func NewRedisConsumer(client redis.UniversalClient) *redisConsumer {
	return &redisConsumer{Client: client}
}

// 接收消息并处理
func (c *redisConsumer) ReceiveMessages(ctx context.Context) {
	for str, p := range MessageProcessor {
		go func(messageKey string, p processor) {
			pubsub := c.Client.Subscribe(ctx, messageKey)
			defer pubsub.Close()
			for {
				ch := pubsub.Channel()
				for msg := range ch {
					if err := p(ctx, msg.Payload); err != nil {
						logger.Error("message processor error")
					}
				}
			}
		}(str, p)
	}
}
