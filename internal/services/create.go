package services

import (
	"context"

	"github.com/neracastle/go-libs/pkg/sys/logger"
	"golang.org/x/exp/slog"

	"github.com/neracastle/chat-server/internal/domain/chat"
)

func (s *Service) Create(ctx context.Context) (int64, error) {
	log := logger.GetLogger(ctx)
	ch := chat.NewChat()
	log.Debug("creating chat")
	err := s.chatRepository.Save(ctx, &ch)

	if err != nil {
		log.Error("failed to create chat", slog.String("error", err.Error()))
	}
	log.Debug("saved chat", ch)
	return ch.Id, err
}
