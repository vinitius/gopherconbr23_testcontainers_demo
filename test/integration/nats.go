package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"

	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

type NatsServer struct {
	instance testcontainers.Container
}

func NewNatsServer(t *testing.T) *NatsServer {
	t.Helper()
	timeout := 3 * time.Minute
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()
	req := testcontainers.ContainerRequest{
		Image:        "nats:2-alpine",
		ExposedPorts: []string{"4222/tcp", "8222/tcp"},
		Cmd:          []string{"--http_port", "8222"},
		WaitingFor:   wait.ForHTTP("/healthz").WithPort("8222"),
		HostConfigModifier: func(config *container.HostConfig) {
			config.AutoRemove = true
		},
	}
	jetStream, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	return &NatsServer{
		instance: jetStream,
	}
}

func (js *NatsServer) Port(t *testing.T) int {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()
	p, err := js.instance.MappedPort(ctx, "4222")
	require.NoError(t, err)
	return p.Int()
}

func (js *NatsServer) Address(t *testing.T) string {
	return fmt.Sprintf("nats://127.0.0.1:%d", js.Port(t))
}

func (js *NatsServer) Close(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()
	require.NoError(t, js.instance.Terminate(ctx))
}
