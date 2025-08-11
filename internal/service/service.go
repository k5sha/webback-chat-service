package service

import (
	"context"
	"fmt"
	"log"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/k5sha/webback-chat-service/internal/mapper"
	"github.com/k5sha/webback-chat-service/internal/store"
	pbAuth "github.com/k5sha/webback-go-proto/gen/go/protos/auth"
	pb "github.com/k5sha/webback-go-proto/gen/go/protos/chat"
)

func (h *Handler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	msg := &store.Message{
		SenderId: req.Sender,
		Content:  req.Content,
		RoomID:   req.RoomId,
		Type:     req.Type,
	}

	message, err := h.store.Messages.Create(ctx, msg)
	if err != nil {
		log.Println("Error creating message:", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	chatMessage := mapper.ToMessageProto(message)
	fmt.Println(chatMessage)
	return &pb.SendMessageResponse{Message: chatMessage}, nil
}

func (h *Handler) GetMessageHistory(ctx context.Context, req *pb.GetMessageHistoryRequest) (*pb.GetMessageHistoryResponse, error) {
	fmt.Println("GetMessageHistory called with request:", req)
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	amount, messages, err := h.store.Messages.ListByRoomID(ctx, req.RoomId, int64(req.Limit), int64(req.Offset))
	if err != nil {
		fmt.Println("Error fetching messages:", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, msg := range messages {
		if msg.SenderId != "" {
			res, err := h.authClient.GetUsername(ctx, &pbAuth.GetUsernameRequest{UserId: msg.SenderId})
			if err != nil {
				log.Printf("failed to get username for user %s: %v", msg.SenderId, err)
				msg.SenderUsername = "Deleted account"
			} else {
				msg.SenderUsername = res.Username
			}
		}

	}

	return &pb.GetMessageHistoryResponse{
		Messages: mapper.ToMessagesProto(messages),
		Amount:   amount,
	}, nil
}
