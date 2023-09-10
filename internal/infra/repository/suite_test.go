package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"go.uber.org/goleak"

	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"gopherconbr.org/23/testcontainers/demo/test/integration"
)

const (
	_initSQLPath = "../../../test/testdata/sql"
)

type SQLRepositorySuite struct {
	suite.Suite
	Ctx                       context.Context
	DB                        bun.IDB
	sqlConn                   *sql.DB
	PostgresDatabaseContainer *integration.PostgresDatabase
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func (s *SQLRepositorySuite) SetupSuite() {
	ctx, cancelFn := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancelFn()

	s.PostgresDatabaseContainer = integration.NewPostgresDatabase(s.T(), _initSQLPath)
	s.sqlConn = sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(s.PostgresDatabaseContainer.DSN(s.T()))))
	s.Require().NoError(s.sqlConn.PingContext(ctx))
	s.DB = bun.NewDB(s.sqlConn, pgdialect.New())
}

func (s *SQLRepositorySuite) TearDownSuite() {
	_ = s.sqlConn.Close()
	s.PostgresDatabaseContainer.Close(s.T())
}
