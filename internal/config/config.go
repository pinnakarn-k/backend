package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Service string
	Env     string
	Port    int
}

func Load() Config {
	loadDotEnv()

	cfg := Config{
		Service: env("SERVICE", "backend"),
		Env:     env("ENV", "local"),
		Port:    envInt("PORT", 8080),
	}

	cfg.validate()

	return cfg
}

func loadDotEnv() {
	_ = godotenv.Load()
}

func (c Config) validate() {
	if c.Service == "" {
		panic("SERVICE is required")
	}

	if c.Port <= 0 {
		panic("PORT must be greater than 0")
	}
}

func (c Config) ListenAddress() string {
	return fmt.Sprintf(":%d", c.Port)
}

func env(key string, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return def
}

func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		panic(fmt.Sprintf("%s must be int", key))
	}

	return i
}
