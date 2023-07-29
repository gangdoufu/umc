package log

import (
	"context"
	"encoding/json"
	"github.com/gangdoufu/umc/pkg/common"
	mycontext "github.com/gangdoufu/umc/pkg/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func defaultOptions() *Option {
	return &Option{
		Path:         "/var/log/umc/test",
		MaxAge:       1,
		MaxSize:      1,
		Level:        "debug",
		LogInConsole: true,
		Format:       "json",
		LevelEncode:  "lower",
	}
}

// 测试配置的等级转换为zap的等级
func TestOptionLevel(t *testing.T) {
	op := defaultOptions()
	levelMap := map[string]zapcore.Level{
		"Info":   zapcore.InfoLevel,
		"DeBug":  zapcore.DebugLevel,
		"ERROR":  zapcore.ErrorLevel,
		"WARN":   zapcore.WarnLevel,
		"DPANIC": zapcore.DPanicLevel,
		"Panic":  zapcore.PanicLevel,
		"FATAL":  zapcore.FatalLevel,
		"OTHER":  zapcore.DebugLevel,
		"":       zapcore.DebugLevel,
	}
	for s, level := range levelMap {
		op.Level = s
		assert.Equal(t, op.zapLevel(), level)
	}
}

func TestOptionEncoderLevel(t *testing.T) {
	op := defaultOptions()
	encoderLevelMap := map[string]zapcore.LevelEncoder{
		"lower":         zapcore.LowercaseLevelEncoder,
		"lower_color":   zapcore.LowercaseColorLevelEncoder,
		"capital":       zapcore.CapitalLevelEncoder,
		"capital_color": zapcore.CapitalColorLevelEncoder,
		"":              zapcore.LowercaseLevelEncoder,
	}
	for s, level := range encoderLevelMap {
		op.LevelEncode = s
		lev := op.zapEncodeLevel()
		assert.IsType(t, lev, level)
	}
}

func TestBuild(t *testing.T) {
	op := defaultOptions()
	cleanPath(op.Path)
	logger := Build(op)
	logger.Info("test", zap.String("test", "test"))
	tempMap := getLogJsonInfo(op.Path)
	assert.Equal(t, tempMap["test"], "test")
}

func cleanPath(path string) {
	if common.FileExist(path) {
		os.RemoveAll(path)
	}
}

func TestGetLogger(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	str := uuid.New().String()
	op := defaultOptions()
	cleanPath(op.Path)
	Build(op)
	infoCtx := mycontext.WithUserInfo(mycontext.WithRequestInfo(ctx, str, now), 1, "testaccount")
	logger := GetLogger(infoCtx)
	logger.Info("测试")
	tempMap := getLogJsonInfo(op.Path)
	assert.Equal(t, tempMap[requestKey], str)
	assert.Equal(t, tempMap[account], "testaccount")
}

func getLogJsonInfo(path string) map[string]string {
	var tempMap map[string]string
	dir, err := os.ReadDir(path)
	if err != nil {
		return tempMap
	}
	var fileInfo []byte
	for _, entry := range dir {
		if !entry.IsDir() {
			info, _ := entry.Info()
			fileInfo, _ = os.ReadFile(filepath.Join(path, info.Name()))
		}
	}

	json.Unmarshal(fileInfo, &tempMap)
	return tempMap
}
func getLogFileInfo(path string) []byte {
	var fileInfo []byte
	dir, err := os.ReadDir(path)
	if err != nil {
		return fileInfo
	}
	for _, entry := range dir {
		if !entry.IsDir() {
			info, _ := entry.Info()
			fileInfo, _ = os.ReadFile(filepath.Join(path, info.Name()))
		}
	}
	return fileInfo
}

func TestGetLoggerWithText(t *testing.T) {
	op := defaultOptions()
	op.Format = "text"
	cleanPath(op.Path)
	Build(op)
	logger := GetLogger(context.Background())
	logger.Info("test", zap.String("test_key", "test_value"))
	info := getLogFileInfo(op.Path)
	assert.True(t, len(info) > 0)
}

func TestContextWithLogger(t *testing.T) {
	op := defaultOptions()
	logger := Build(op)
	ctx := ContextWithLogger(context.Background(), logger)
	contextLogger := GetContextLogger(ctx)
	assert.Equal(t, logger, contextLogger)
}
