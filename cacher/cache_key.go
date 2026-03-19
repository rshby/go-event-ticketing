package cacher

import (
	"fmt"

	"github.com/rshby/go-event-ticketing/config"
)

func GetEventByIDCacheKey(id uint64) string {
	return createCacheKey(fmt.Sprintf("event:id:%d", id))
}

func createCacheKey(key string) string {
	return fmt.Sprintf("ticketing:%s:%s", config.Mode(), key)
}
