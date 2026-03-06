package cacher

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rshby/go-event-ticketing/config"
	"github.com/rshby/go-event-ticketing/tracing"
)

type CacheManager interface {
	Set(ctx context.Context, key string, item *Item) error
}

type cacheManager struct {
	client *redis.Client
}

// NewCacheManager creates new instance of CacheManager
func NewCacheManager(client *redis.Client) CacheManager {
	return &cacheManager{
		client: client,
	}
}

func (c *cacheManager) Set(ctx context.Context, key string, item *Item) error {
	ctx, span := tracing.Start(ctx)
	defer span.End()

	// check if no ttl
	var err error
	if item.IsNoTTL() {
		err = c.client.Set(ctx, key, item.value, 0).Err()
	} else {
		var ttl = item.TTL()
		if ttl == 0 {
			ttl = config.DefaultRedisTTL
		}

		err = c.client.Set(ctx, key, item.value, ttl).Err()
	}

	// check if any error
	if err != nil {
		return err
	}

	return nil
}
