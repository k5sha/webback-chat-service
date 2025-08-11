package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/k5sha/webback-chat-service/internal/config"
	"github.com/k5sha/webback-chat-service/internal/db"
	"github.com/k5sha/webback-chat-service/internal/service"
	"github.com/k5sha/webback-chat-service/internal/store"
	pbAuth "github.com/k5sha/webback-go-proto/gen/go/protos/auth"
	pb "github.com/k5sha/webback-go-proto/gen/go/protos/chat"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()

	// gRPC Listener
	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}
	defer lis.Close()

	// Database
	db, err := db.New(cfg.DB.Addr, cfg.DB.MaxOpenConns, cfg.DB.MaxIdleConns, cfg.DB.MaxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	log.Println("database connection pool established")

	// Migration

	if _, err := os.Stat(cfg.DB.MigrationsPath); os.IsNotExist(err) {
		log.Fatalf("Migrations path does not exist: %s", cfg.DB.MigrationsPath)
	}

	sourceURL := "file://" + cfg.DB.MigrationsPath

	log.Println("Starting DB migrations")
	log.Printf("Using migrations path: %s", cfg.DB.MigrationsPath)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %v", err)
	}
	log.Printf("Current working dir: %s", dir)

	m, err := migrate.New(sourceURL, cfg.DB.Addr)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	err = m.Up()
	if err != nil && migrate.ErrNoChange != err {
		log.Fatalf("Migration failed: %v", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No new migrations to run")
	} else {
		log.Println("Migrations applied successfully")
	}

	// Store
	store := store.NewStorage(db)

	// gRPC connection to auth service
	authConn, err := grpc.NewClient(cfg.AuthServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to auth service: %v", err)
	}
	defer authConn.Close()

	log.Println("successfuly connect to auth service:", cfg.AuthServiceAddr)

	// Grpc init
	grpcAuthClient := pbAuth.NewAuthServiceClient(authConn)

	chatHandler := service.NewHandler(cfg, store, grpcAuthClient)

	// gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, chatHandler)

	// Run gRPC server in goroutine
	go func() {
		log.Println("gRPC server started at", cfg.GRPCAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP gateway (gRPC-Gateway)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(gwmux)

	if err := pb.RegisterChatServiceHandlerFromEndpoint(ctx, gwmux, cfg.GRPCAddr, opts); err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	log.Println("HTTP gateway started at", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, handler); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
