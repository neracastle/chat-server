package grpc_server

import (
	"golang.org/x/exp/slog"

	"github.com/neracastle/chat-server/internal/services"
	userdesc "github.com/neracastle/chat-server/pkg/chat_v1"
)

type Server struct {
	userdesc.UnimplementedChatV1Server
	logger      *slog.Logger
	chatService *services.Service
}

func (s *Server) GetLogger() *slog.Logger {
	return s.logger
}

func NewServer(logger *slog.Logger, srv *services.Service) *Server {
	return &Server{logger: logger, chatService: srv}
}
