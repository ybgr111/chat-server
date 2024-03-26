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

type createChatVariables struct {
	id        int64
	chatModel model.Chat
}

type CreateChatSuite struct {
	ctx       context.Context
	ctxWithTx context.Context

	suite.Suite

	mc                 *minimock.Controller
	chatRepositoryMock *repoMocks.ChatRepositoryMock
	fakeTxMock         *mocks.FakeTxMock
	transactorMock     *dbMocks.TransactorMock

	txManagerMock db.TxManager

	createChatVariables
}

func TestCreateChatSuite(t *testing.T) {
	suite.Run(t, new(CreateChatSuite))
}

func (s *CreateChatSuite) SetupSuite() {
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
	s.chatModel = model.Chat{
		Usernames: []string{gofakeit.Username(), gofakeit.Username(), gofakeit.Username()},
	}
}

func (s *CreateChatSuite) TestCreateChat_Success() {
	// Специфичные моки методов.
	s.chatRepositoryMock.CreateMock.Expect(s.ctxWithTx, converter.ToChatCreate(&s.chatModel)).Return(s.id, nil)
	s.fakeTxMock.CommitMock.Return(nil)

	service := chat.NewService(s.chatRepositoryMock, s.txManagerMock)

	newID, err := service.Create(s.ctx, &s.chatModel)

	// Проверки корректности теста.
	require.Nil(s.T(), nil, err)
	require.Equal(s.T(), s.id, newID)
}

func (s *CreateChatSuite) TestCreateChat_FailCreateChat() {
	chatErr := errors.New("cant create chat")

	s.chatRepositoryMock.CreateMock.Return(0, chatErr)
	s.fakeTxMock.RollbackMock.Return(nil)

	service := chat.NewService(s.chatRepositoryMock, s.txManagerMock)

	_, err := service.Create(s.ctx, &s.chatModel)

	require.Error(s.T(), chatErr, err)
}
