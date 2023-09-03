package stream

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type AccountsSubscriber struct {
	natsConn *nats.Conn
}

func NewAccountsSubscriber(natsConn *nats.Conn) AccountsSubscriber {
	return AccountsSubscriber{natsConn: natsConn}
}

func (s AccountsSubscriber) Subscribe(ctx context.Context, rcv chan<- []byte, ready chan<- struct{}) error {
	const (
		fetchTimeout = 200 * time.Millisecond
	)
	defer func() {
		err := s.natsConn.Drain()
		if err != nil {
			if !errors.Is(err, nats.ErrConnectionReconnecting) { // No point in draining if connection was lost
				log.Fatalf("failed to drain nats connection: %v", err)
			}
		}
	}()

	sub, err := s.natsConn.SubscribeSync("ACCOUNTS.new")
	if err != nil {
		return err
	}
	ready <- struct{}{}
	fmt.Println("==subscription ready")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("==subscription ended")
			return nil
		default:
		}

		msg, err := sub.NextMsg(fetchTimeout)
		if err != nil {
			if errors.Is(err, nats.ErrTimeout) {
				continue
			}
			return err
		}

		fmt.Printf("==received msg: %q on subject %q", string(msg.Data), msg.Subject)
		rcv <- msg.Data
	}
}
