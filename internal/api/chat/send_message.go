package chat

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ybgr111/chat-server/internal/converter"
	desc "github.com/ybgr111/chat-server/pkg/note_v1"
)

func (i *Server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*empty.Empty, error) {
	err := i.chatService.SendMessage(ctx, converter.ToMessage(req.GetMessage()))
	if err != nil {
		return nil, err
	}

	log.Print("message delivered\n")

	return &empty.Empty{}, nil
}
