package chat

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	desc "github.com/ybgr111/chat-server/pkg/note_v1"
)

func (i *Server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	err := i.chatService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("chat with id: %d deleted\n", req.GetId())

	return &empty.Empty{}, nil
}
