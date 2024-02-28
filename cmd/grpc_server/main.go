package main

import (
	"context"
	"flag"
	"log"
	"net"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ybgr111/chat-server/internal/config"
	"github.com/ybgr111/chat-server/internal/config/env"
	desc "github.com/ybgr111/chat-server/pkg/note_v1"
)

var configPath string

const (
	chatTable       = "chat"
	messageTable    = "message"
	idColumn        = "id"
	usernamesColumn = "usernames"
	fromColumn      = `"from"`
	textColumn      = "text"
	timestampColumn = "timestamp"
)

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedNoteV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	builderInsert := sq.Insert(chatTable).
		PlaceholderFormat(sq.Dollar).
		Columns(usernamesColumn).
		Values(req.Users.Usernames).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to build query")
	}

	var authID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&authID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to insert chat")
	}

	log.Printf("Chat members: %v", req.GetUsers())

	return &desc.CreateResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	builderUpdate := sq.Delete(chatTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.Id})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to build query")
	}

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to delete chat")
	}

	if res.RowsAffected() == 0 {
		return nil, errors.WithMessage(errors.New("failed to delete chat"), "chat not found")
	}

	log.Printf("Chat id: %d", req.GetId())
	return &empty.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*empty.Empty, error) {
	builderInsert := sq.Insert(messageTable).
		PlaceholderFormat(sq.Dollar).
		Columns(fromColumn, textColumn, timestampColumn).
		Values(req.Message.From, req.Message.Text, req.Message.Timestamp.AsTime()).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to build query")
	}

	var authID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&authID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to insert message")
	}

	log.Printf("Chat message: %v", req.GetMessage())

	return &empty.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: #{err}")
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: #{err}")
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
