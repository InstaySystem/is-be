package initialization

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/InstaySystem/is-be/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

func InitRedis(cfg *config.Config) (*redis.Client, error) {
	rAddr := cfg.Redis.Host + fmt.Sprintf(":%d", cfg.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:      rAddr,
		Password:  cfg.Redis.Password,
		TLSConfig: &tls.Config{},
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("cache - %w", err)
	}

	return rdb, nil
}
