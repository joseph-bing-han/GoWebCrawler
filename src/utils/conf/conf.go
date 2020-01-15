package conf

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Get(key string, defaultValue string) string {

	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
