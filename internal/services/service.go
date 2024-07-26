package services

import (
	"github.com/neracastle/chat-server/internal/repository"
)

type Service struct {
	chatRepository repository.Repository
}

func NewService(chatRepository repository.Repository) *Service {
	return &Service{chatRepository: chatRepository}
}
