package main

import (
	"context"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/ybgr111/chat-server/pkg/note_v1"
)

const (
	address = "localhost:50053"
	chatID  = 12
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	create, err := c.Create(ctx, &desc.CreateRequest{
		Users: &desc.Users{
			Usernames: []string{"Nikita", "Eshe Nikita"},
		},
	})
	if err != nil {
		log.Fatalf("failed to create chat: %v", err)
	}
	log.Printf(color.RedString("Сreate chat info:\n"), color.GreenString("%+v", create.GetId()))

	delete, err := c.Delete(ctx, &desc.DeleteRequest{
		Id: chatID,
	})
	if err != nil {
		log.Fatalf("failed to delete chat by id: %v", err)
	}
	log.Printf(color.RedString("Delete chat info:\n"), color.GreenString("%+v", delete.String()))

	sendMessage, err := c.SendMessage(ctx, &desc.SendMessageRequest{
		Message: &desc.Message{
			From:      gofakeit.Name(),
			Text:      gofakeit.BeerHop(),
			Timestamp: timestamppb.New(gofakeit.Date()),
		},
	})
	if err != nil {
		log.Fatalf("failed to send massage: %v", err)
	}
	log.Printf(color.RedString("Send message info:\n"), color.GreenString("%+v", sendMessage.String()))
}
