package jwt

import (
	"context"
	"fmt"
	"time"

	rdb "github.com/redis/go-redis/v9"
)

const blacklistPrefix = "jwt:blacklist:"

type Blacklist struct {
	client *rdb.Client
}

func NewBlacklist(client *rdb.Client) *Blacklist {
	return &Blacklist{client: client}
}

// Add 将 token 加入黑名单，TTL 与 token 剩余有效期一致
func (b *Blacklist) Add(ctx context.Context, tokenStr string, expiry time.Time) error {
	ttl := time.Until(expiry)
	if ttl <= 0 {
		return nil // 已过期，无需加入
	}
	key := blacklistPrefix + tokenStr
	return b.client.Set(ctx, key, 1, ttl).Err()
}

// IsBlocked 返回 true 表示 token 已被拉黑
func (b *Blacklist) IsBlocked(ctx context.Context, tokenStr string) (bool, error) {
	key := blacklistPrefix + tokenStr
	exists, err := b.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("blacklist check: %w", err)
	}
	return exists > 0, nil
}
