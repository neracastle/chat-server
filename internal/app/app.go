package app

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/neracastle/go-libs/pkg/closer"
	"github.com/neracastle/go-libs/pkg/sys/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpc_server "github.com/neracastle/chat-server/internal/grpc-server"
	"github.com/neracastle/chat-server/pkg/chat_v1"
)

type App struct {
	grpc        *grpc.Server
	srvProvider *serviceProvider
}

func NewApp(ctx context.Context) *App {
	app := &App{srvProvider: newServiceProvider()}
	app.init(ctx)
	return app
}

func (a *App) init(ctx context.Context) {
	lg := logger.SetupLogger(a.srvProvider.Config().Env)
	a.grpc = grpc.NewServer()

	reflection.Register(a.grpc)
	chat_v1.RegisterChatV1Server(a.grpc, grpc_server.NewServer(lg, a.srvProvider.ChatService(ctx)))
}

func (a *App) Start() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	conn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.srvProvider.Config().GRPC.Host, a.srvProvider.Config().GRPC.Port))
	if err != nil {
		return err
	}

	log.Printf("ChatAPI service started on %s:%d\n", a.srvProvider.Config().GRPC.Host, a.srvProvider.Config().GRPC.Port)

	closer.Add(func() error {
		a.grpc.GracefulStop()
		return nil
	})

	if err = a.grpc.Serve(conn); err != nil {
		return err
	}

	return nil
}
