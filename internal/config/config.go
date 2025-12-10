// Package config
package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Address        string
	AllowedOrigins []string
	Interval       time.Duration
	LogLevel       string
	LogFormat      string
	JWTSecret      string
	JWTExpiry      time.Duration
	DatabaseURL    string
}

func Load() *Config {
	_ = godotenv.Load()

	addr := getEnv("HTTP_ADDR", ":3000")

	var origins []string
	rawOrigins := os.Getenv("ALLOWED_ORIGINS")
	if rawOrigins != "" {
		parts := strings.SplitSeq(rawOrigins, ",")
		for o := range parts {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				origins = append(origins, trimmed)
			}
		}
	}

	interval := 3 * time.Second
	if raw := os.Getenv("SCRAPE_INTERVAL"); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil && parsed > 0 {
			interval = parsed
		}
	}

	jwtExpiry := 24 * time.Hour
	if raw := os.Getenv("JWT_EXPIRY"); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil && parsed > 0 {
			jwtExpiry = parsed
		}
	}

	return &Config{
		Address:        addr,
		AllowedOrigins: origins,
		Interval:       interval,
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		LogFormat:      getEnv("LOG_FORMAT", "text"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpiry:      jwtExpiry,
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/horizonx"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
