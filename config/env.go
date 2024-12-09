package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string
	DBName     string
	DBUser     string
	DBPassword string
	// DBAddress string
	TokenSymmetricKey   string
	AccessTokenDuration time.Duration
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func InitConfig() Config {
	godotenv.Load()
	// Parse the TOKEN_DURATION string into a time.Duration
	duration, err := time.ParseDuration(getEnv("TOKEN_DURATION", "24h"))
	if err != nil {
		// If there's an error, use the default duration "24h"
		duration = 24 * time.Hour
	}

	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port:       getEnv("PORT", "8080"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		// DBAddress: fmt.Sprintf("%s:%s",getEnv("DB_HOST","localhost"),getEnv("DB_PORT","3306")),
		DBName:              getEnv("DB_NAME", "simple_bank"),
		TokenSymmetricKey:   getEnv("TOKEN_SYMMETRIC_KEY", ""),
		AccessTokenDuration: duration,
	}
}

var Envs = InitConfig()
