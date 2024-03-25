package chat

import (
	"context"

	"github.com/ybgr111/chat-server/internal/model"
	"github.com/ybgr111/chat-server/internal/repository/chat/converter"
)

func (s *serv) Create(ctx context.Context, chat *model.Chat) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.chatRepository.Create(ctx, converter.ToChatCreate(chat))
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
