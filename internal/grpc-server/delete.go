package grpc_server

import (
	"context"

	"github.com/neracastle/go-libs/pkg/sys/logger"
	"google.golang.org/protobuf/types/known/emptypb"

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
