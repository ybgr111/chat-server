package converter

import (
	"github.com/ybgr111/chat-server/internal/model"
	modelRepo "github.com/ybgr111/chat-server/internal/repository/chat/model"
)

func ToChatCreate(
	chat *model.Chat,
) *modelRepo.Chat {
	return &modelRepo.Chat{
		Usernames: chat.Usernames,
	}
}

func ToSendMessage(
	message *model.Message,
) *modelRepo.Message {
	return &modelRepo.Message{
		From: message.From,
		Text: message.Text,
	}
}
