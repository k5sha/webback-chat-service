package service

import (
	"context"
	"log"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/k5sha/webback-chat-service/internal/mapper"
	"github.com/k5sha/webback-chat-service/internal/store"
	pb "github.com/k5sha/webback-go-proto/gen/go/protos/chat"
)

func (h *Handler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	log.Println(req)
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	log.Println(req)
	msg := &store.Message{
		Sender:  req.Sender,
		Content: req.Content,
		RoomID:  req.RoomId,
		Type:    req.Type,
	}

	message, err := h.store.Messages.Create(ctx, msg)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	chatMessage := mapper.ToMessageProto(message)

	return &pb.SendMessageResponse{Message: chatMessage}, nil
}

func (h *Handler) GetMessageHistory(ctx context.Context, req *pb.GetMessageHistoryRequest) (*pb.GetMessageHistoryResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	messages, err := h.store.Messages.ListByRoomID(ctx, req.RoomId, int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetMessageHistoryResponse{
		Messages: mapper.ToMessagesProto(messages),
	}, nil
}
