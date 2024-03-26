package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/ybgr111/platform_common/pkg/db"
	dbMocks "github.com/ybgr111/platform_common/pkg/db/mocks"
	"github.com/ybgr111/platform_common/pkg/db/pg"
	"github.com/ybgr111/platform_common/pkg/db/transaction"

	"github.com/ybgr111/chat-server/internal/model"
	"github.com/ybgr111/chat-server/internal/repository/chat/converter"
	repoMocks "github.com/ybgr111/chat-server/internal/repository/mocks"
	"github.com/ybgr111/chat-server/internal/service/chat"
	"github.com/ybgr111/chat-server/internal/service/mocks"
)

type sendMessageVariables struct {
	id           int64
	messageModel model.Message
}

type SendMessageSuite struct {
	ctx       context.Context
	ctxWithTx context.Context

	suite.Suite

	mc                 *minimock.Controller
	chatRepositoryMock *repoMocks.ChatRepositoryMock
	fakeTxMock         *mocks.FakeTxMock
	transactorMock     *dbMocks.TransactorMock

	txManagerMock db.TxManager

	sendMessageVariables
}

func TestSendMessageSuite(t *testing.T) {
	suite.Run(t, new(SendMessageSuite))
}

func (s *SendMessageSuite) SetupSuite() {
	s.ctx = context.Background()
	s.mc = minimock.NewController(s.T())

	s.chatRepositoryMock = repoMocks.NewChatRepositoryMock(s.mc)
	s.fakeTxMock = mocks.NewFakeTxMock(s.mc)

	s.ctxWithTx = pg.MakeContextTx(s.ctx, s.fakeTxMock)
	s.transactorMock = dbMocks.NewTransactorMock(s.mc)
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	s.transactorMock.BeginTxMock.Expect(s.ctx, txOpts).Return(s.fakeTxMock, nil)

	s.txManagerMock = transaction.NewTransactionManager(s.transactorMock)

	s.id = gofakeit.Int64()
	s.messageModel = model.Message{
		From: gofakeit.Username(),
		Text: gofakeit.Sentence(10),
	}
}

func (s *SendMessageSuite) TestSendMessage_Success() {
	// Специфичные моки методов.
	s.chatRepositoryMock.SendMessageMock.Expect(s.ctxWithTx, converter.ToSendMessage(&s.messageModel)).Return(nil)
	s.fakeTxMock.CommitMock.Return(nil)

	service := chat.NewService(s.chatRepositoryMock, s.txManagerMock)

	err := service.SendMessage(s.ctx, &s.messageModel)

	// Проверки корректности теста.
	require.Nil(s.T(), nil, err)
}

func (s *SendMessageSuite) TestSendMessage_FailSendMessage() {
	messageErr := errors.New("cant send message")

	s.chatRepositoryMock.SendMessageMock.Return(messageErr)
	s.fakeTxMock.RollbackMock.Return(nil)

	service := chat.NewService(s.chatRepositoryMock, s.txManagerMock)

	err := service.SendMessage(s.ctx, &s.messageModel)

	require.Error(s.T(), messageErr, err)
}
