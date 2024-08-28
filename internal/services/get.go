package services

import (
	"context"
	"errors"

	syserr "github.com/neracastle/go-libs/pkg/sys/error"

	"github.com/neracastle/chat-server/internal/domain/chat"
	"github.com/neracastle/chat-server/internal/repository"
)

func (s *Service) Get(ctx context.Context, chatId int64) (*chat.Chat, error) {
	ch, err := s.chatRepository.Get(ctx, chatId)
	if err != nil {
		if errors.Is(err, repository.ErrChatNotFound) {
			return nil, syserr.New("Чат не найден", syserr.NotFound)
		}

		return nil, err
	}

	return ch, nil
}
