package service

import (
	"github.com/k5sha/webback-chat-service/internal/config"
	"github.com/k5sha/webback-chat-service/internal/store"
	pb "github.com/k5sha/webback-go-proto/gen/go/protos/chat"
)

type Handler struct {
	pb.UnimplementedChatServiceServer
	config config.Config
	store  store.Storage
}

func NewHandler(
	config config.Config,
	store store.Storage) *Handler {
	return &Handler{
		config: config,
		store:  store,
	}
}
