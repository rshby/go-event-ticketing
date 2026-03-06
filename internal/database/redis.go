package database

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rshby/go-event-ticketing/config"
	"github.com/sirupsen/logrus"
)

// ConnectRedis connects redis
func ConnectRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", config.RedisHost(), config.RedisPort()),
		DB:              config.RedisDbNumber(),
		PoolSize:        config.RedisMaxConnSize(),
		MinIdleConns:    config.RedisIdleConnSize(),
		ConnMaxLifetime: config.RedisConnLifetime(),
	})

	logrus.Infof("success connect Redis [%s:%s db number %d]✅", config.RedisHost(), config.RedisPort(), config.RedisDbNumber())
	return redisClient
}
