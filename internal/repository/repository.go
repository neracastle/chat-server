package repository

import (
	"context"

	domain "github.com/neracastle/chat-server/internal/domain/chat"
)

type Repository interface {
	Save(context.Context, *domain.Chat) error
	Delete(ctx context.Context, id int32) error
}
