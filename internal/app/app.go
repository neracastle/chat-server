package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpc_server "github.com/neracastle/chat-server/internal/grpc-server"
	"github.com/neracastle/chat-server/pkg/chat_v1"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
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
	lg := setupLogger(a.srvProvider.Config().Env)
	a.grpc = grpc.NewServer()

	reflection.Register(a.grpc)
	chat_v1.RegisterChatV1Server(a.grpc, grpc_server.NewServer(lg, a.srvProvider.ChatService(ctx)))
}

func (a *App) Start() error {

	conn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.srvProvider.Config().GRPC.Host, a.srvProvider.Config().GRPC.Port))
	if err != nil {
		return err
	}

	log.Printf("ChatAPI service started on %s:%d\n", a.srvProvider.Config().GRPC.Host, a.srvProvider.Config().GRPC.Port)

	if err = a.grpc.Serve(conn); err != nil {
		return err
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) {

	allClosed := make(chan struct{})
	go func() {
		_ = a.srvProvider.dbc.Close()
		a.grpc.GracefulStop()
		close(allClosed)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-allClosed:
			return
		}
	}
}

func setupLogger(env string) *slog.Logger {
	var lg *slog.Logger

	switch env {
	case envLocal:
		lg = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		lg = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		lg = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return lg
}
