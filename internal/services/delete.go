package services

import (
	"context"

	"github.com/neracastle/go-libs/pkg/sys/logger"
	"golang.org/x/exp/slog"
)

func (s *Service) Delete(ctx context.Context, chatId int64) error {
	log := logger.GetLogger(ctx)
	err := s.chatRepository.Delete(ctx, chatId)

	if err != nil {
		log.Error("failed to delete chat", slog.String("error", err.Error()), slog.Int("chatId", int(chatId)))
	}

	return err
}
