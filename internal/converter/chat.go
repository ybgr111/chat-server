package converter

import (
	"github.com/ybgr111/chat-server/internal/model"
	desc "github.com/ybgr111/chat-server/pkg/note_v1"
)

func ToChat(users *desc.Users) *model.Chat {
	return &model.Chat{
		Usernames: users.Usernames,
	}
}

func ToMessage(message *desc.Message) *model.Message {
	return &model.Message{
		From: message.From,
		Text: message.Text,
	}
}
