package chat

import (
	"github.com/ybgr111/chat-server/internal/client/db"
	"github.com/ybgr111/chat-server/internal/repository"
	"github.com/ybgr111/chat-server/internal/service"
)

type serv struct {
	chatRepository repository.ChatRepository
	txManager      db.TxManager
}

func NewService(
	chatRepository repository.ChatRepository,
	txManager db.TxManager,
) service.ChatService {
	return &serv{
		chatRepository: chatRepository,
		txManager:      txManager,
	}
}
