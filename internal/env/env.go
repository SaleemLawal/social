package env

import (
	"os"
	"strconv"
)

func GetString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetInt(key string, fallback int) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return val
}

func GetBool(key string, fallback bool) bool {
	val, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return val
}
