package grpc_server

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"golang.org/x/exp/slog"

	userdesc "github.com/neracastle/chat-server/pkg/chat_v1"
)

// Create создает чат
func (s *Server) Create(ctx context.Context, req *userdesc.CreateRequest) (*userdesc.CreateResponse, error) {
	log := s.GetLogger()
	log = log.With(slog.String("method", "grpc-server.Create"))

	log.Debug("called", slog.Any("users", req.Usernames))

	var id int64
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

	err = tx.QueryRow(ctx, "INSERT INTO chat.chats(created_at) VALUES (now()) RETURNING id").Scan(&id)
	if err != nil {
		log.Error("failed to save chat in db", slog.String("error", err.Error()))
		return nil, err
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	insertQuery := psql.Insert("chat.chat_users").Columns("chat_id", "user_id", "user_name")

	for idx, username := range req.Usernames {
		insertQuery = insertQuery.Values(id, idx, username)
	}

	query, args, err := insertQuery.ToSql()
	if err != nil {
		log.Error("failed to create insert query", slog.String("error", err.Error()))
		return nil, err
	}

	res, err := tx.Exec(ctx, query, args...)
	if err != nil {
		log.Error("failed to save chat in db", slog.String("error", err.Error()))
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error("failed to save chat in db", slog.String("error", err.Error()))
		return nil, err
	}

	log.Debug("saved chat in db",
		slog.Int64("id", id),
		slog.Int64("affected", res.RowsAffected()))

	return &userdesc.CreateResponse{Id: id}, nil
}
