package log

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type myGormLogger struct {
	logger        *zap.Logger
	level         gormlogger.LogLevel
	SlowThreshold time.Duration // 慢sql
}

func NewGormLogger(logger *zap.Logger) *myGormLogger {
	return &myGormLogger{logger: logger}
}

func (l *myGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.level = level
	return &newLogger
}

func (l myGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogger.Info {
		GetLoggerWithContext(ctx, l.logger).Sugar().Infof(msg, data...)
	}
}

func (l myGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogger.Warn {
		GetLoggerWithContext(ctx, l.logger).Sugar().Warnf(msg, data...)
	}
}
func (l myGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogger.Error {
		GetLoggerWithContext(ctx, l.logger).Sugar().Errorf(msg, data...)
	}
}

func (l myGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	logger := GetLoggerWithContext(ctx, l.logger)
	var curLevel gormlogger.LogLevel

	if (err != nil && !errors.Is(err, gorm.ErrRecordNotFound)) || l.level >= gormlogger.Error {
		curLevel = gormlogger.Error
	} else if (l.SlowThreshold >= 0 && elapsed >= l.SlowThreshold) || l.level >= gormlogger.Warn {
		curLevel = gormlogger.Warn
	} else if l.level >= gormlogger.Info {
		curLevel = gormlogger.Info
	}
	sql, rows := fc()
	var fields []zap.Field
	fields = append(fields, zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	switch curLevel {
	case gormlogger.Error:
		fields = append(fields, zap.Error(err))
		logger.Error("trace", fields...)
	case gormlogger.Warn:
		logger.Warn("trace", fields...)
	case gormlogger.Info:
		logger.Info("trace", fields...)
	}
}

var (
	gormPackage = filepath.Join("gorm.io", "gorm")
)

func (l myGormLogger) getLogger(ctx context.Context) *zap.Logger {
	logger := GetLoggerWithContext(ctx, l.logger)
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		default:
			return logger.WithOptions(zap.AddCallerSkip(i))
		}
	}
	return logger
}
