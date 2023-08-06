package redis

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const (
	tokenBlacklistKey = "token_black_list"
	userJwtToken      = "user_jwt_token"
	tenantAppTenant   = "tenant_appid_token"
)

// 将 token 加入到黑名单
func AddTokenTOBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	key := strings.Join([]string{tokenBlacklistKey, token}, defaultStep)
	return client.Set(ctx, key, "token", expiration).Err()
}

// 检查token 是否在黑名单中
func CheckTokenInBlack(ctx context.Context, token string) bool {
	key := strings.Join([]string{tokenBlacklistKey, token}, defaultStep)
	s := client.Exists(ctx, key).Val()
	val := client.Get(ctx, key).Val()
	fmt.Print(val)
	return s > 0
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

// 每次取用都延期,
func CheckTenantDynamicToken(ctx context.Context, appId string, token string) bool {
	key := strings.Join([]string{tenantAppTenant, appId, token}, defaultStep)
	exists := client.Exists(ctx, key).Val() > 0
	if exists {
		client.PExpire(ctx, key, 20*time.Minute)
	}
	return client.Exists(ctx, key).Val() > 0
}

// 租户的动态token 有效期只有20分钟 主要用于接口调试,过期了就要重新搞
func SetTenantDynamicToken(ctx context.Context, appId string, token string) error {
	key := strings.Join([]string{tenantAppTenant, appId, token}, defaultStep)
	return client.Set(ctx, key, "token", 20*time.Minute).Err()
}
