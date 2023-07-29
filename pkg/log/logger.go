package log

import (
	"context"
	"github.com/gangdoufu/umc/pkg/common"
	mycontext "github.com/gangdoufu/umc/pkg/context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger *zap.Logger

func newLogger(op *Option) {
	encoder := op.newEncoder()
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= op.zapLevel() && level < zap.ErrorLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	})
	var cores []zapcore.Core
	cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(op.getWriter("info.log")), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(op.getWriter("error.log")), errorLevel))
	// 如果要在控制台输出,需要增加控制台输出
	if op.LogInConsole {
		allLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level >= op.zapLevel()
		})
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), allLevel))
	}
	core := zapcore.NewTee(cores...)
	logger = zap.New(core, zap.AddCaller())
}

const (
	requestKey = "request-id"
	requestAt  = "request-at"
	account    = "account"
)

// GetLogger 获取定制化的日志
func GetLogger(ctx context.Context) *zap.Logger {
	return logger.With(getContextInfoFields(ctx)...)
}

func Logger() *zap.Logger {
	return logger
}

func getContextInfoFields(ctx context.Context) []zap.Field {
	requestInfo := mycontext.GetRequestInfo(ctx)
	user := mycontext.GetContextUserInfo(ctx)
	var fields []zap.Field
	if requestInfo != nil {
		fields = append(fields, zap.String(requestKey, requestInfo.RequestId), zap.Time(requestAt, requestInfo.RequestAt))
	}
	if user != nil {
		fields = append(fields, zap.String(account, user.Account))
	}
	return fields
}

func GetLoggerWithContext(ctx context.Context, l *zap.Logger) *zap.Logger {
	return l.With(getContextInfoFields(ctx)...)
}

func ContextWithLogger(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

func GetContextLogger(ctx context.Context) *zap.Logger {
	return common.GetContextValue[zap.Logger](ctx, loggerKey)
}
