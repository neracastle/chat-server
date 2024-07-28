package services

import (
	"context"

	"golang.org/x/exp/slog"

	"github.com/neracastle/chat-server/internal/app/logger"
	"github.com/neracastle/chat-server/internal/domain/chat"
	"github.com/neracastle/chat-server/internal/services/models"
)

func (s *Service) Create(ctx context.Context, req models.Create) (int32, error) {
	log := logger.GetLogger(ctx)
	ch := chat.NewChat(req.UserIds)
	err := s.chatRepository.Save(ctx, &ch)

	if err != nil {
		log.Error("failed to create chat", slog.String("error", err.Error()), slog.Any("request", req))
	}

	return ch.Id, err
}
