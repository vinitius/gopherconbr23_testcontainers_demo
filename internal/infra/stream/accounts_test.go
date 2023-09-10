package stream

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
)

type AccountsStreamSuite struct {
	NatsSuite
}

func TestAccountsStreamSuite(t *testing.T) {
	suite.Run(t, new(AccountsStreamSuite))
}

func (s *AccountsStreamSuite) SetupTest() {
	s.Ctx = context.TODO()
}

/*
TestSubscribe garante que um AccountsSubscriber recebe todas as mensagens postadas no tópico ACCOUNTS.new
assim que uma inscrição é realizada com sucesso.
*/
func (s *AccountsStreamSuite) TestSubscribe() {
	// Arrange
	const (
		timeout = 3 * time.Second
	)
	var (
		expected = [][]byte{
			[]byte(`{"account_id":"2", "type":"regular", "label":"misc", "currency":"BTC"}`),
			[]byte(`{"account_id":"3", "type":"regular", "label":"misc", "currency":"BRL"}`),
		}
		actual [][]byte
	)

	ctx, cancelFn := context.WithTimeout(s.Ctx, timeout)
	defer cancelFn()
	underTest := NewAccountsSubscriber(s.NatsConn)
	received := make(chan []byte, len(expected))
	ready := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case data, ok := <-received:
				if ok {
					actual = append(actual, data)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Act
	go func() {
		err := underTest.Subscribe(ctx, received, ready)
		s.Require().NoError(err)
	}()
	<-ready
	close(ready)
	for _, msg := range expected {
		err := s.NatsConn.Publish("ACCOUNTS.new", msg)
		s.Require().NoError(err)
	}

	// Assert
	s.Require().Eventually(func() bool {
		return len(actual) == len(expected)
	}, timeout, 200*time.Millisecond, "did not receive expected messages in time")
	s.Require().Empty(cmp.Diff(expected, actual))
}
