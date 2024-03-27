package chat

import (
	"github.com/ybgr111/chat-server/internal/service"
	desc "github.com/ybgr111/chat-server/pkg/chat_v1"
)

type Server struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService
}

func NewServer(userService service.ChatService) *Server {
	return &Server{
		chatService: userService,
	}
}
