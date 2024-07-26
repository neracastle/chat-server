package app

import (
	"context"
	"log"

	db "github.com/neracastle/chat-server/internal/client"
	"github.com/neracastle/chat-server/internal/client/pg"
	"github.com/neracastle/chat-server/internal/config"
	"github.com/neracastle/chat-server/internal/repository"
	"github.com/neracastle/chat-server/internal/repository/postgres"
	"github.com/neracastle/chat-server/internal/services"
)

type serviceProvider struct {
	conf           *config.Config
	chatService    *services.Service
	chatRepository repository.Repository
	dbc            db.Client
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) Config() config.Config {
	if sp.conf == nil {
		cfg := config.MustLoad()
		sp.conf = &cfg
	}

	return *sp.conf
}

func (sp *serviceProvider) DbClient(ctx context.Context) db.Client {
	if sp.dbc == nil {
		client, err := pg.NewClient(ctx, sp.Config().Postgres.DSN())
		if err != nil {
			log.Fatalf("failed to connect to pg: %v", err)
		}

		err = client.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("failed ping to pg: %v", err)
		}

		sp.dbc = client
	}

	return sp.dbc
}

func (sp *serviceProvider) ChatRepository(ctx context.Context) repository.Repository {
	if sp.chatRepository == nil {
		sp.chatRepository = postgres.New(sp.DbClient(ctx))
	}

	return sp.chatRepository
}

func (sp *serviceProvider) ChatService(ctx context.Context) *services.Service {
	if sp.chatService == nil {
		sp.chatService = services.NewService(sp.ChatRepository(ctx))
	}

	return sp.chatService
}
