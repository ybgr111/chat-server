package chat

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	desc "github.com/ybgr111/chat-server/pkg/chat_v1"
)

func (i *Server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	err := i.chatService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
