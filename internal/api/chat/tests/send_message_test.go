package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"

	"github.com/ybgr111/chat-server/internal/api/chat"
	"github.com/ybgr111/chat-server/internal/model"
	"github.com/ybgr111/chat-server/internal/service"
	serviceMocks "github.com/ybgr111/chat-server/internal/service/mocks"
	desc "github.com/ybgr111/chat-server/pkg/chat_v1"
)

func TestSendMessage(t *testing.T) {
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.SendMessageRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		from = gofakeit.Username()
		text = gofakeit.Sentence(10)

		serviceErr = fmt.Errorf("service error")

		req = &desc.SendMessageRequest{
			Message: &desc.Message{
				From: from,
				Text: text,
			},
		}

		messageModel = &model.Message{
			From: from,
			Text: text,
		}

		res = &empty.Empty{}
	)

	tests := []struct {
		name            string
		args            args
		want            *empty.Empty
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, messageModel).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, messageModel).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {

			chatServiceMock := tt.chatServiceMock(mc)
			api := chat.NewServer(chatServiceMock)

			newID, err := api.SendMessage(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
