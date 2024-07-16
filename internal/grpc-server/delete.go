package grpc_server

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"golang.org/x/exp/slog"
	"google.golang.org/protobuf/types/known/emptypb"

	userdesc "github.com/neracastle/chat-server/pkg/chat_v1"
)

// Delete удаляет чат
func (s *Server) Delete(ctx context.Context, req *userdesc.DeleteRequest) (*emptypb.Empty, error) {
	log := s.GetLogger()
	log = log.With(slog.String("method", "grpc-server.Delete"))

	log.Debug("called", slog.Int64("id", req.GetId()))

	tx, err := s.pgcon.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction chat in db", slog.String("error", err.Error()))
		return nil, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if !errors.Is(err, pgx.ErrTxClosed) {
			log.Error("failed to rollback transaction", slog.String("error", err.Error()))
		}
	}(tx, ctx)

	res, err := tx.Exec(ctx, "DELETE FROM chat.chat_users WHERE chat_id=$1", req.GetId())
	if err != nil {
		log.Error("failed to delete chat", slog.String("error", err.Error()))
		return nil, err
	}

	res, err = tx.Exec(ctx, "DELETE FROM chat.chats WHERE id=$1", req.GetId())
	if err != nil {
		log.Error("failed to delete chat", slog.String("error", err.Error()))
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error("failed to delete chat in db", slog.String("error", err.Error()))
		return nil, err
	}

	log.Debug("deleted chat",
		slog.Int64("id", req.GetId()),
		slog.Int64("rows", res.RowsAffected()))

	return &emptypb.Empty{}, nil
}
