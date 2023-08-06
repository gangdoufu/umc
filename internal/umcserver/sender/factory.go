package sender

import "context"

func GetSender(sType senderType) ISender {
	switch sType {
	case email:
		return emailSender
	default:
		return emailSender
	}
}

type ISender interface {
	SendMessage(ctx context.Context, vo *InfoVo) error
}
