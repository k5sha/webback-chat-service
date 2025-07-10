package config

import (
	"time"

	"github.com/k5sha/webback-chat-service/internal/env"
)

type Config struct {
	Addr string
	DB   dbConfig
}

type dbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

func Load() Config {

	return Config{
		Addr: env.GetString("ADDR", ":3003"),
		DB: dbConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5433/chat?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetDuration("DB_MAX_IDLE_TIME", 5*time.Minute),
		},
	}
}
