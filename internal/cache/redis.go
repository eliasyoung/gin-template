package cache

import (
	"context"
	"fmt"

	"github.com/eliasyoung/gin-template/internal/config"
	"github.com/eliasyoung/gin-template/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func InitRedis(cfg *config.CacheConfig) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.Get().Info("redis_addr", zap.String("addr", addr))

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("could not ping redis: %w", err)
	}

	return client, nil
}
