package grpc_server

import (
	"context"

	"github.com/neracastle/chat-server/internal/app/logger"
	"github.com/neracastle/chat-server/internal/services/models"
	userdesc "github.com/neracastle/chat-server/pkg/chat_v1"
)

// Create создает чат
func (s *Server) Create(ctx context.Context, req *userdesc.CreateRequest) (*userdesc.CreateResponse, error) {
	ctx = logger.AssignLogger(ctx, s.GetLogger())
	id, err := s.chatService.Create(ctx, models.Create{UserIds: req.GetUserIds()})

	if err != nil {
		return nil, err
	}

	return &userdesc.CreateResponse{Id: id}, nil
}
