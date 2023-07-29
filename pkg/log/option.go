package log

import (
	"github.com/gangdoufu/umc/pkg/common"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"path/filepath"
	"strings"
	"time"
)

const (
	infoLevel   = "info"
	debugLevel  = "debug"
	errorLevel  = "error"
	warnLevel   = "warn"
	dpanicLevel = "dpanic"
	panicLevel  = "panic"
	fatalLevel  = "fatal"

	loggerKey = "my_logger"
)

type Option struct {
	Path         string `mapstructure:"path" json:"path" yaml:"path" `                               // 日志路径
	MaxAge       int    `mapstructure:"max-age" json:"max-age" yaml:"max-age" `                      // 日志最大保存的天数
	MaxSize      int64  `mapstructure:"max-size" json:"max-size" yaml:"max-size" `                   // 日志切割体积
	Level        string `mapstructure:"level" json:"level" yaml:"level" `                            // 记录的起始等级
	LogInConsole bool   `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console" ` // 是否在控制台输出
	Format       string `mapstructure:"format" json:"format" yaml:"format" `
	LevelEncode  string `mapstructure:"level-encode" json:"level-encode" yaml:"level-encode" `
}

// 获取配置的 日志等级显示格式 小写，小写带颜色 大写 大写带颜色
func (o Option) zapEncodeLevel() zapcore.LevelEncoder {
	switch strings.ToLower(o.LevelEncode) {
	case "lower":
		return zapcore.LowercaseLevelEncoder
	case "lower_color":
		return zapcore.LowercaseColorLevelEncoder
	case "capital":
		return zapcore.CapitalLevelEncoder
	case "capital_color":
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

// 根据配置中的信息获取zap对于的等级
func (o Option) zapLevel() zapcore.Level {
	level := strings.ToLower(o.Level)
	switch level {
	case infoLevel:
		return zapcore.InfoLevel
	case debugLevel:
		return zapcore.DebugLevel
	case errorLevel:
		return zapcore.ErrorLevel
	case warnLevel:
		return zapcore.WarnLevel
	case dpanicLevel:
		return zapcore.DPanicLevel
	case panicLevel:
		return zapcore.PanicLevel
	case fatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

func Build(op *Option) *zap.Logger {
	newLogger(op)
	return logger
}

// 编码器配置
func (o Option) newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		EncodeLevel:    o.zapEncodeLevel(),
		EncodeTime:     o.timeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// 日期的格式
func (o Option) timeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// 获取日志写入的writer 使用 rotatelog 工具包对日志文件进行切割
func (o Option) getWriter(filename string) io.Writer {
	file := filepath.Join(o.Path, filename)
	hook, err := rotatelogs.New(file+".%Y%m%d%H",
		rotatelogs.WithMaxAge(time.Duration(o.MaxAge*24)*time.Hour),
		rotatelogs.WithRotationSize(o.MaxSize*common.MB))
	if err != nil {
		panic(err)
	}
	return hook
}

// 创建编码器,目前支持两种,
func (o Option) newEncoder() zapcore.Encoder {
	if o.Format == "json" {
		return zapcore.NewJSONEncoder(o.newEncoderConfig())
	} else {
		return zapcore.NewConsoleEncoder(o.newEncoderConfig())
	}
}
