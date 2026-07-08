package configs

import "os"

func getEnv(k, defaultValue string) string {
	if value := os.Getenv(k); value != "" {
		return value
	}

	return defaultValue
}
