package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	}
}

func GetEnv(key string) string {
	fmt.Println(key)
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("❌ Required environment variable %s is not set", key)
	}
	return value
}
