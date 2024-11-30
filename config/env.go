package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string
	DBName     string
	DBUser     string
	DBPassword string
	// DBAddress string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port:       getEnv("PORT", "8080"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		// DBAddress: fmt.Sprintf("%s:%s",getEnv("DB_HOST","localhost"),getEnv("DB_PORT","3306")),
		DBName: getEnv("DB_NAME", "simple_bank"),
	}
}

var Envs = initConfig()
