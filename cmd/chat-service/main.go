package main

import (
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/k5sha/webback-chat-service/internal/config"
	"github.com/k5sha/webback-chat-service/internal/db"
	"github.com/k5sha/webback-chat-service/internal/service"
	"github.com/k5sha/webback-chat-service/internal/store"
	pb "github.com/k5sha/webback-go-proto/gen/go/protos/chat"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		log.Fatal("failed to listen:", cfg.Addr)
	}
	defer l.Close()

	// Database
	db, err := db.New(cfg.DB.Addr, cfg.DB.MaxOpenConns, cfg.DB.MaxIdleConns, cfg.DB.MaxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	log.Println("database connection pool established")

	// Migration

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %v", err)
	}

	migrationsPath := filepath.Join(dir, "cmd", "migrate", "migrations")
	sourceURL := "file://" + filepath.ToSlash(migrationsPath)

	m, err := migrate.New(
		sourceURL,
		cfg.DB.Addr,
	)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run up migrations: %v", err)
	}

	// Store
	store := store.NewStorage(db)

	chatHandler := service.NewHandler(cfg, store)

	pb.RegisterChatServiceServer(grpcServer, chatHandler)

	log.Println("GPRC server has started at", cfg.Addr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
