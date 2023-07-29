package context

import (
	"context"
	"github.com/gangdoufu/umc/pkg/common"
	"gorm.io/gorm"
	"time"
)

const (
	userInfoKey        = "my_user_info"
	requestInfoKey     = "my_request_info"
	transactionInfoKey = "my_transaction_info"
)

func WithUserInfo(ctx context.Context, userId uint, account string) context.Context {
	user := &UserVo{UserId: userId, Account: account}
	return context.WithValue(ctx, userInfoKey, user)
}

func WithRequestInfo(ctx context.Context, requestId string, requestAt time.Time) context.Context {
	request := &RequestVo{RequestId: requestId, RequestAt: requestAt}
	return context.WithValue(ctx, requestInfoKey, request)
}

func WithTransaction(ctx context.Context, db *gorm.DB) context.Context {
	t := &TransactionVo{Db: db}
	return context.WithValue(ctx, transactionInfoKey, t)
}

func GetContextUserInfo(ctx context.Context) *UserVo {
	return common.GetContextValue[UserVo](ctx, userInfoKey)

}
func GetRequestInfo(ctx context.Context) *RequestVo {
	return common.GetContextValue[RequestVo](ctx, requestInfoKey)
}
func GetTransactionInfo(ctx context.Context) *TransactionVo {
	return common.GetContextValue[TransactionVo](ctx, transactionInfoKey)
}
