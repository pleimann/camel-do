package utils

import (
	"os"
)

func EnvWithDefault(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)

	if !exists || len(value) == 0 {
		return defaultValue
	}

	return value
}
