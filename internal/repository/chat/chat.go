package chat

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/ybgr111/chat-server/internal/repository"
	chatModel "github.com/ybgr111/chat-server/internal/repository/chat/model"
	"github.com/ybgr111/platform_common/pkg/db"
)

const (
	chatTable       = "chat"
	messageTable    = "message"
	idColumn        = "id"
	usernamesColumn = "usernames"
	fromColumn      = `"from"`
	textColumn      = "text"
	timestampColumn = "timestamp"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, chat *chatModel.Chat) (int64, error) {
	builderInsert := sq.Insert(chatTable).
		PlaceholderFormat(sq.Dollar).
		Columns(usernamesColumn).
		Values(chat.Usernames).
		Suffix(fmt.Sprintf("RETURNING %s", idColumn))

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "chat_repository.Create",
		QueryRaw: query,
	}

	var chatID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)
	if err != nil {
		return 0, errors.WithMessage(err, "failed to insert chat")
	}

	return chatID, nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builderUpdate := sq.Delete(chatTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return errors.WithMessage(err, "failed to build query")
	}

	q := db.Query{
		Name:     "chat_repository.Delete",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return errors.WithMessage(err, "failed to delete chat")
	}

	if res.RowsAffected() == 0 {
		return errors.WithMessage(errors.New("failed to delete chat"), "chat not found")
	}

	return nil
}

func (r *repo) SendMessage(ctx context.Context, message *chatModel.Message) error {
	builderInsert := sq.Insert(messageTable).
		PlaceholderFormat(sq.Dollar).
		Columns(fromColumn, textColumn, timestampColumn).
		Values(message.From, message.Text, time.Now()).
		Suffix(fmt.Sprintf("RETURNING %s", idColumn))

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_repository.SendMessage",
		QueryRaw: query,
	}

	var messageID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&messageID)
	if err != nil {
		return errors.WithMessage(err, "failed to insert message")
	}

	return nil
}
