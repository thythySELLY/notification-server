package config

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis/v7"
)

var RedisClient *redis.Client

func InitRedis() {
	host := GetEnv("REDIS_HOST")
	port := GetEnv("REDIS_PORT")
	password := GetEnv("REDIS_PASSWORD")
	dbStr := GetEnv("REDIS_DB")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		log.Fatalf("❌ Error converting REDIS_DB: %v", err)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	_, err = RedisClient.Ping().Result()
	if err != nil {
		log.Fatalf("❌ Cannot connect to Redis: %v", err)
	}

	fmt.Println("✅ Redis connection successful!")
}
