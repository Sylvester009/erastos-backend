package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL"`
	Port        string `env:"PORT"`
	JWTSecret   string `env:"JWT_SECRET"`
	CORSOrigins []string
}

func LoadConfig() (*Config, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		CORSOrigins: corsOrigins(os.Getenv("CORS_ORIGINS")),
	}

	return config, nil
}

func corsOrigins(value string) []string {
	if strings.TrimSpace(value) == "" {
		return []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}

	var origins []string
	for _, origin := range strings.Split(value, ",") {
		if origin = strings.TrimSpace(origin); origin != "" {
			origins = append(origins, origin)
		}
	}
	if len(origins) == 0 {
		return []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}
	return origins
}
