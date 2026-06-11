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

	// SQL Server
	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string

	// Redis
	RedisEnabled  bool
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int
}

func Load() Config {
	loadDotEnv()

	cfg := Config{
		Service: env("SERVICE", "backend"),
		Env:     env("ENV", "local"),
		Port:    envInt("PORT", 8080),

		DBHost:     env("DB_HOST", ""),
		DBPort:     envInt("DB_PORT", 1433),
		DBName:     env("DB_NAME", ""),
		DBUser:     env("DB_USER", ""),
		DBPassword: env("DB_PASSWORD", ""),

		RedisEnabled:  envBool("REDIS_ENABLED", false),
		RedisHost:     env("REDIS_HOST", "localhost"),
		RedisPort:     envInt("REDIS_PORT", 6379),
		RedisPassword: env("REDIS_PASSWORD", ""),
		RedisDB:       envInt("REDIS_DB", 0),
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

func envBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		panic(fmt.Sprintf("%s must be bool", key))
	}

	return b
}
