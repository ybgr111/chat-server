package chat

import (
	"context"

	"github.com/ybgr111/chat-server/internal/model"
	"github.com/ybgr111/chat-server/internal/repository/chat/converter"
)

func (s *serv) SendMessage(ctx context.Context, message *model.Message) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.chatRepository.SendMessage(ctx, converter.ToSendMessage(message))
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
