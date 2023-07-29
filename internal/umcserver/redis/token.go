package redis

import (
	"context"
	"strings"
	"time"
)

const (
	tokenBlacklistKey = "token_black_list"
	userJwtToken      = "user_jwt_token"
)

// 将 token 加入到黑名单
func AddTokenTOBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	key := strings.Join([]string{tokenBlacklistKey, token}, defaultStep)
	return client.Set(ctx, key, "token", expiration).Err()
}

// 检查token 是否在黑名单中
func CheckTokenInBlack(ctx context.Context, token string) bool {
	key := strings.Join([]string{tokenBlacklistKey, token}, defaultStep)
	s := client.Get(ctx, key).String()
	return s == "token"
}

func AddUserToken(ctx context.Context, token string, userId string, expiration time.Duration) error {
	key := strings.Join([]string{userJwtToken, token}, defaultStep)
	return client.Set(ctx, key, userId, expiration).Err()
}

func CheckUserToken(ctx context.Context, token string, userId string) bool {
	key := strings.Join([]string{userJwtToken, token}, defaultStep)
	u := client.Get(ctx, key).String()
	return userId == u
}
