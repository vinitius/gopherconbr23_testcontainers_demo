package app

import (
	"context"
	"encoding/json"

	"golang.org/x/sync/errgroup"

	"gopherconbr.org/23/testcontainers/demo/internal/infra/repository"
	"gopherconbr.org/23/testcontainers/demo/internal/infra/stream"
)

type AccountsHandlerFn func(ctx context.Context, ch <-chan []byte) error

type Application struct {
	subscriber stream.AccountsSubscriber
	repository repository.AccountsSQLRepository
	handler    AccountsHandlerFn
	ready      chan struct{}
}

func NewApplication(subscriber stream.AccountsSubscriber,
	repo repository.AccountsSQLRepository,
	ready chan struct{},
) Application {
	return Application{
		subscriber: subscriber,
		repository: repo,
		ready:      ready,
		handler: func(ctx context.Context, ch <-chan []byte) error {
			for {
				select {
				case data, ok := <-ch:
					if ok {
						var account repository.AccountModel
						if err := json.Unmarshal(data, &account); err != nil {
							return err
						}
						if err := repo.Create(ctx, account); err != nil {
							return err
						}
					}
				case <-ctx.Done():
					return nil
				}
			}
		},
	}
}

func (a Application) Run(ctx context.Context) error {
	errGroup, ctx := errgroup.WithContext(ctx)
	received := make(chan []byte)

	errGroup.Go(func() error {
		return a.handler(ctx, received)
	})

	errGroup.Go(func() error {
		return a.subscriber.Subscribe(ctx, received, a.ready)
	})

	return errGroup.Wait()
}
