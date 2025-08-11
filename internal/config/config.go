package config

import (
	"time"

	"github.com/k5sha/webback-chat-service/internal/env"
)

type Config struct {
	Addr            string
	GRPCAddr        string
	AuthServiceAddr string
	DB              dbConfig
}

type dbConfig struct {
	Addr           string
	MaxOpenConns   int
	MaxIdleConns   int
	MaxIdleTime    time.Duration
	MigrationsPath string
}

func Load() Config {

	return Config{
		Addr:            env.GetString("ADDR", ":8083"),
		GRPCAddr:        env.GetString("GRPC_ADDR", ":3003"),
		AuthServiceAddr: env.GetString("AUTH_SERVICE_ADDR", ":3004"),
		DB: dbConfig{
			Addr:           env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/chat?sslmode=disable"),
			MaxOpenConns:   env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns:   env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:    env.GetDuration("DB_MAX_IDLE_TIME", 5*time.Minute),
			MigrationsPath: env.GetString("DB_MIGRATIONS_PATH", "./cmd/migrate/migrations"),
		},
	}
}
