// Package config
package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Address string
}

func Load() *Config {
	godotenv.Load()

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":3000"
	}
	return &Config{Address: addr}
}
