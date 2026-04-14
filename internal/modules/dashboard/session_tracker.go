package dashboard

import (
	"context"
	"strconv"
	"time"

	rdb "github.com/redis/go-redis/v9"
)

type SessionTracker interface {
	MarkOnline(ctx context.Context, userID uint, ttl time.Duration) error
}

type redisSessionTracker struct {
	client *rdb.Client
}

func NewSessionTracker(client *rdb.Client) SessionTracker {
	if client == nil {
		return noopSessionTracker{}
	}
	return &redisSessionTracker{client: client}
}

func (t *redisSessionTracker) MarkOnline(ctx context.Context, userID uint, ttl time.Duration) error {
	if ttl <= 0 {
		return nil
	}
	key := onlineUserKeyPrefix + strconv.FormatUint(uint64(userID), 10)
	return t.client.Set(ctx, key, time.Now().UTC().Format(time.RFC3339), ttl).Err()
}

type noopSessionTracker struct{}

func (noopSessionTracker) MarkOnline(context.Context, uint, time.Duration) error {
	return nil
}
