package redis

import (
	"context"
	"go-server/internal/config"

	rdb "github.com/redis/go-redis/v9"
)

func Init(cfg config.RedisConfig) (*rdb.Client, error) {
	client := rdb.NewClient(&rdb.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
