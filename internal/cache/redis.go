package cache

import (
	"banner-system/internal/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func OpenCache(cacheOpts config.CacheOpts) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cacheOpts.Host, cacheOpts.Port),
		Password: cacheOpts.Password,
		DB:       cacheOpts.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("error on ping Redis server, %s", err)
	}

	return client, nil
}
