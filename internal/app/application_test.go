package app

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.uber.org/goleak"

	"gopherconbr.org/23/testcontainers/demo/internal/infra/repository"
	"gopherconbr.org/23/testcontainers/demo/internal/infra/stream"
	"gopherconbr.org/23/testcontainers/demo/test/integration"
)

const (
	_initSQLPath = "../../test/testdata/sql"
)

type AppSuite struct {
	suite.Suite
	Ctx                       context.Context
	DB                        bun.IDB
	sqlConn                   *sql.DB
	NatsConn                  *nats.Conn
	NatsContainer             *integration.NatsServer
	PostgresDatabaseContainer *integration.PostgresDatabase
	subscriber                stream.AccountsSubscriber
	repository                repository.AccountsSQLRepository
	ready                     chan struct{}
}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func (s *AppSuite) SetupSuite() {
	ctx, cancelFn := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancelFn()

	s.PostgresDatabaseContainer = integration.NewPostgresDatabase(s.T(), _initSQLPath)
	s.sqlConn = sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(s.PostgresDatabaseContainer.DSN(s.T()))))
	s.Require().NoError(s.sqlConn.PingContext(ctx))
	s.DB = bun.NewDB(s.sqlConn, pgdialect.New())
	s.repository = repository.NewAccountsSQLRepository(s.DB)

	s.NatsContainer = integration.NewNatsServer(s.T())
	natsConn, err := nats.Connect(s.NatsContainer.Address(s.T()))
	s.Require().NoError(err)
	s.NatsConn = natsConn
	s.subscriber = stream.NewAccountsSubscriber(s.NatsConn)

	s.ready = make(chan struct{}, 1)
}

func (s *AppSuite) TearDownSuite() {
	_ = s.sqlConn.Close()
	s.PostgresDatabaseContainer.Close(s.T())
	s.NatsContainer.Close(s.T())
}

func (s *AppSuite) SetupTest() {
	s.Ctx = context.TODO()

	// Truncate before every test
	_, err := s.DB.NewTruncateTable().Model(new(repository.AccountModel)).Exec(s.Ctx)
	s.Require().NoError(err)

	s.Require().NoError(err)
}

/*
TestRun garante que a aplicação inicia, se conecta ao tópico ACCOUNTS.new e persiste os dados dos eventos recebidos
na tabela 'accounts'.
*/
func (s *AppSuite) TestRun() {
	// Arrange
	timeout := 3 * time.Second
	ctx, cancelFn := context.WithTimeout(s.Ctx, timeout)
	defer cancelFn()
	underTest := NewApplication(s.subscriber, s.repository, s.ready)
	expected := &repository.AccountModel{
		ID:       "0",
		Type:     "regular",
		Label:    "e2e",
		Currency: "BRL",
	}

	// Act
	go func() {
		err := underTest.Run(ctx)
		s.Require().NoError(err)
	}()

	<-s.ready
	err := s.NatsConn.Publish("ACCOUNTS.new", []byte(`{"account_id":"0", "type":"regular", "label":"e2e", "currency":"BRL"}`))
	s.Require().NoError(err)

	// Assert
	s.Require().Eventually(func() bool {
		exists, err := s.DB.NewSelect().Model(expected).Exists(ctx)
		s.Require().NoError(err)
		return exists
	}, timeout, 200*time.Millisecond)

	actual := &repository.AccountModel{}
	err = s.DB.NewSelect().Model(expected).WherePK().Scan(ctx, actual)
	s.Require().NoError(err)
	s.Require().Empty(cmp.Diff(expected, actual))
}
