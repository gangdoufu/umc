package common

import (
	"context"
	"crypto/md5"
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cast"
	"os"
	"strings"
)

func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func GetContextValue[T any](ctx context.Context, key string) *T {
	value := ctx.Value(key)
	if value == nil {
		return nil
	}
	return getValue[T](value)
}

func getValue[T any](value interface{}) *T {
	if vo, ok := value.(*T); !ok {
		return nil
	} else {
		return vo
	}
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func IntSliceToString(ints []int) string {
	var res strings.Builder
	for _, u := range ints {
		res.WriteString(cast.ToString(u))
	}
	return res.String()
}

func RandToken(num int) string {
	b := make([]byte, num)
	crand.Read(b)
	return fmt.Sprintf("%x", b)
}
