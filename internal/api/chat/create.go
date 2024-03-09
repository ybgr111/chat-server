package chat

import (
	"context"
	"log"

	"github.com/ybgr111/chat-server/internal/converter"
	desc "github.com/ybgr111/chat-server/pkg/note_v1"
)

func (i *Server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.chatService.Create(ctx, converter.ToChat(req.GetUsers()))

	if err != nil {
		return nil, err
	}

	log.Printf("id: %d\n", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
