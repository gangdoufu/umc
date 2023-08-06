package consumer

import (
	"context"
	"encoding/json"
	"github.com/gangdoufu/umc/internal/umcserver/message"
	"github.com/gangdoufu/umc/internal/umcserver/sender"
)

// 接收到消息之后的处理方法
type processor func(ctx context.Context, str string) error

// 所有需要处理的消息放在这里,key 是消息名称, value 是消息处理器
var MessageProcessor = map[string]processor{
	message.VerificationCodeKey: SendVerificationCode,
}

type Consumer interface {
	ReceiveMessages(ctx context.Context)
}

func SendVerificationCode(ctx context.Context, info string) error {
	var infoVo = &sender.InfoVo{}
	err := json.Unmarshal([]byte(info), infoVo)
	if err != nil {
		return err
	}
	curSender := sender.GetSender(infoVo.Type)
	return curSender.SendMessage(ctx, infoVo)

}
