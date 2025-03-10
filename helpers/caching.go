package helpers

import (
	"notification-server/config"
	"time"
)

func SetCache(key string, value string) error {
	expiration := 1 * time.Minute
	return config.RedisClient.Set(key, value, expiration).Err()
}

func GetCache(key string) (string, error) {
	return config.RedisClient.Get(key).Result()
}

func DeleteCache(key string) error {
	return config.RedisClient.Del(key).Err()
}
