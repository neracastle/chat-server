package grpc_server

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/neracastle/chat-server/internal/app/logger"
	userdesc "github.com/neracastle/chat-server/pkg/chat_v1"
)

// Delete удаляет чат
func (s *Server) Delete(ctx context.Context, req *userdesc.DeleteRequest) (*emptypb.Empty, error) {
	ctx = logger.AssignLogger(ctx, s.GetLogger())
	err := s.chatService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
