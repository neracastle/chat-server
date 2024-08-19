package services

import (
	"context"

	"github.com/neracastle/chat-server/internal/repository"
	"github.com/neracastle/chat-server/internal/services/models"
)

type ChatService interface {
	Create(ctx context.Context, req models.Create) (int64, error)
	Delete(ctx context.Context, chatId int64) error
}

type AuthServiceClient interface {
	CanDelete(ctx context.Context, userID int64) (bool, error)
}

type Service struct {
	chatRepository repository.Repository
	authService    AuthServiceClient
}

func NewService(chatRepository repository.Repository, authClient AuthServiceClient) *Service {
	return &Service{
		chatRepository: chatRepository,
		authService:    authClient,
	}
}
