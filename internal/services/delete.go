package services

import (
	"context"

	"golang.org/x/exp/slog"

	"github.com/neracastle/chat-server/internal/app/logger"
)

func (s *Service) Delete(ctx context.Context, chatId int32) error {
	log := logger.GetLogger(ctx)
	err := s.chatRepository.Delete(ctx, chatId)

	if err != nil {
		log.Error("failed to delete chat", slog.String("error", err.Error()), slog.Int("chatId", int(chatId)))
	}

	return err
}
