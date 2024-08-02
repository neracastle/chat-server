package postgres

import (
	"context"
	"fmt"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/neracastle/go-libs/pkg/db"

	domain "github.com/neracastle/chat-server/internal/domain/chat"
	"github.com/neracastle/chat-server/internal/repository"
)

const (
	createdColumn       = "created_at"
	idColumn            = "id"
	usersChatIdColumn   = "chat_id"
	usersUserIdColumn   = "user_id"
	usersUserNameColumn = "user_name"
)

var _ repository.Repository = (*repo)(nil)

type repo struct {
	conn db.Client
}

func New(conn db.Client) repository.Repository {
	instance := &repo{conn: conn}

	return instance
}

func (r *repo) Save(ctx context.Context, chat *domain.Chat) error {
	err := r.conn.DB().ReadCommitted(ctx, func(ctx context.Context) error {
		insert, _, err := sq.Insert("chat.chats").
			Columns(createdColumn).
			Values(sq.Expr("now()")).
			Suffix(fmt.Sprintf("RETURNING %s", idColumn)).
			ToSql()
		if err != nil {
			return err
		}

		q := db.Query{Name: "Save", QueryRaw: insert}

		err = r.conn.DB().QueryRow(ctx, q).Scan(&chat.Id)
		if err != nil {
			return err
		}

		psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
		insertQuery := psql.Insert("chat.chat_users").Columns(usersChatIdColumn, usersUserIdColumn, usersUserNameColumn)

		for idx, userId := range chat.UserIds {
			insertQuery = insertQuery.Values(chat.Id, userId, "U"+strconv.Itoa(idx))
		}

		query, args, err := insertQuery.ToSql()
		if err != nil {
			return err
		}

		q = db.Query{Name: "SaveChatUsers", QueryRaw: query}
		_, err = r.conn.DB().Exec(ctx, q, args...)

		return err
	})

	return err
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	err := r.conn.DB().ReadCommitted(ctx, func(ctx context.Context) error {
		psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
		sql, args, err := psql.Delete("chat.chat_users").
			Where(sq.Eq{usersChatIdColumn: id}).
			ToSql()
		if err != nil {
			return err
		}

		q := db.Query{Name: "DeleteChatUsers", QueryRaw: sql}
		_, err = r.conn.DB().Exec(ctx, q, args...)
		if err != nil {
			return err
		}

		sql, args, err = psql.Delete("chat.chats").
			Where(sq.Eq{idColumn: id}).
			ToSql()
		if err != nil {
			return err
		}

		q = db.Query{Name: "Delete", QueryRaw: sql}
		_, err = r.conn.DB().Exec(ctx, q, args...)

		return err
	})

	return err
}
