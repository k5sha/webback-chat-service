package service

import (
	"github.com/k5sha/webback-chat-service/internal/config"
	"github.com/k5sha/webback-chat-service/internal/store"
	pbAuth "github.com/k5sha/webback-go-proto/gen/go/protos/auth"
	pb "github.com/k5sha/webback-go-proto/gen/go/protos/chat"
)

type Handler struct {
	pb.UnimplementedChatServiceServer
	config     config.Config
	store      store.Storage
	authClient pbAuth.AuthServiceClient
}

func NewHandler(
	config config.Config,
	store store.Storage,
	authClient pbAuth.AuthServiceClient) *Handler {
	return &Handler{
		config:     config,
		store:      store,
		authClient: authClient,
	}
}
