package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port int
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Port: envInt("PORT", 8080),
	}
}

func (c Config) ListenAddress() string {
	return fmt.Sprintf(":%d", c.Port)
}

func envInt(key string, def int) int {
	v := os.Getenv(key)

	if v == "" {
		return def
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		panic(key + " must be integer")
	}

	return i
}
