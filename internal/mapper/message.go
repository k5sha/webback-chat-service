package mapper

import (
	"github.com/k5sha/webback-chat-service/internal/store"
	pb "github.com/k5sha/webback-go-proto/gen/go/protos/chat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToMessageProto(m *store.Message) *pb.ChatMessage {
	return &pb.ChatMessage{
		Id:             m.ID,
		RoomId:         m.RoomID,
		SenderId:       m.SenderId,
		SenderUsername: m.SenderUsername,
		Content:        m.Content,
		Type:           m.Type,
		SentAt:         timestamppb.New(m.SentAt),
	}
}

func ToMessagesProto(messages []*store.Message) []*pb.ChatMessage {
	result := make([]*pb.ChatMessage, len(messages))
	for i, m := range messages {
		result[i] = ToMessageProto(m)
	}
	return result
}
