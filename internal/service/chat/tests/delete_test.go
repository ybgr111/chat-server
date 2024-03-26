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

	repoMocks "github.com/ybgr111/chat-server/internal/repository/mocks"
	"github.com/ybgr111/chat-server/internal/service/chat"
	"github.com/ybgr111/chat-server/internal/service/mocks"
)

type deleteChatVariables struct {
	id int64
}

type DeleteChatSuite struct {
	ctx       context.Context
	ctxWithTx context.Context

	suite.Suite

	mc                 *minimock.Controller
	chatRepositoryMock *repoMocks.ChatRepositoryMock
	fakeTxMock         *mocks.FakeTxMock
	transactorMock     *dbMocks.TransactorMock

	txManagerMock db.TxManager

	deleteChatVariables
}

func TestDeleteChatSuite(t *testing.T) {
	suite.Run(t, new(DeleteChatSuite))
}

func (s *DeleteChatSuite) SetupSuite() {
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
}

func (s *DeleteChatSuite) TestDeleteChat_Success() {
	// Специфичные моки методов.
	s.chatRepositoryMock.DeleteMock.Expect(s.ctxWithTx, s.id).Return(nil)
	s.fakeTxMock.CommitMock.Return(nil)

	service := chat.NewService(s.chatRepositoryMock, s.txManagerMock)

	err := service.Delete(s.ctx, s.id)

	// Проверка корректности теста.
	require.Nil(s.T(), nil, err)
}

func (s *DeleteChatSuite) TestDeleteChat_FailDeleteChat() {
	chatErr := errors.New("cant delete chat")

	s.chatRepositoryMock.DeleteMock.Return(chatErr)
	s.fakeTxMock.RollbackMock.Return(nil)

	service := chat.NewService(s.chatRepositoryMock, s.txManagerMock)

	err := service.Delete(s.ctx, s.id)

	require.Error(s.T(), chatErr, err)
}
