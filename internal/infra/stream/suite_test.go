package stream

import (
	"context"
	"testing"

	"go.uber.org/goleak"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/suite"

	"gopherconbr.org/23/testcontainers/demo/test/integration"
)

type NatsSuite struct {
	suite.Suite
	Ctx           context.Context
	NatsConn      *nats.Conn
	NatsContainer *integration.NatsServer
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func (s *NatsSuite) SetupSuite() {
	s.NatsContainer = integration.NewNatsServer(s.T())
	natsConn, err := nats.Connect(s.NatsContainer.Address(s.T()))
	s.Require().NoError(err)

	s.NatsConn = natsConn
}

func (s *NatsSuite) TearDownSuite() {
	s.NatsContainer.Close(s.T())
}
