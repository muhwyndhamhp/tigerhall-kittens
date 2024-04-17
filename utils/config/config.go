package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	PORT                    = "PORT"
	ENV_FILE                = ".env"
	LIBSQL_URL              = "LIBSQL_URL"
	LIBSQL_TOKEN            = "LIBSQL_TOKEN"
	JWT_SECRET              = "JWT_SECRET"
	JWT_EXPIRY_DURATION     = "JWT_EXPIRY_DURATION"
	CF_ACCOUNT_ID           = "CF_ACCOUNT_ID"
	CF_R2_ACCESS_KEY_ID     = "CF_R2_ACCESS"
	CF_R2_SECRET_ACCESS_KEY = "CF_R2_SECRET"
)

func init() {
	if err := godotenv.Load(ENV_FILE); err != nil {
		fmt.Printf("Failed to load env file: %s\n", err)
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
