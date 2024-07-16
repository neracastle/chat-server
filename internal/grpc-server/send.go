package grpc_server

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	userdesc "github.com/neracastle/chat-server/pkg/chat_v1"
)

// SendMessage отправляет сообщение в чат
func (s *Server) SendMessage(ctx context.Context, req *userdesc.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("called SendMessage method with req: %v", req)

	//пока неясно надо ли хранить сообщения в базе или же просто обмен ими будет идти через кафку

	return &emptypb.Empty{}, nil
}
