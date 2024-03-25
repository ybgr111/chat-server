package repository

import (
	"context"

	chatModel "github.com/ybgr111/chat-server/internal/repository/chat/model"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *chatModel.Chat) (int64, error)
	Delete(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, message *chatModel.Message) error
}
