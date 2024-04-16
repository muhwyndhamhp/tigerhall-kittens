package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	ENV_FILE            = ".env"
	LIBSQL_URL          = "LIBSQL_URL"
	LIBSQL_TOKEN        = "LIBSQL_TOKEN"
	JWT_SECRET          = "JWT_SECRET"
	JWT_EXPIRY_DURATION = "JWT_EXPIRY_DURATION"
)

func init() {
	if err := godotenv.Load(ENV_FILE); err != nil {
		fmt.Printf("Failed to load env file: %s\n", err)
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
