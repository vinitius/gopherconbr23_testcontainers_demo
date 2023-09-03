package integration

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresDatabase struct {
	instance testcontainers.Container
}

func NewPostgresDatabase(t *testing.T, relativeSQLDir string) *PostgresDatabase {
	t.Helper()

	absPath, err := filepath.Abs(relativeSQLDir)
	require.NoError(t, err)

	timeout := 3 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:10",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "accounts_test",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
		Mounts:     testcontainers.Mounts(testcontainers.BindMount(absPath, "/docker-entrypoint-initdb.d")),
		HostConfigModifier: func(config *container.HostConfig) {
			config.AutoRemove = true
		},
	}
	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	return &PostgresDatabase{
		instance: postgres,
	}
}

func (db *PostgresDatabase) Port(t *testing.T) int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	p, err := db.instance.MappedPort(ctx, "5432")
	require.NoError(t, err)
	return p.Int()
}

func (db *PostgresDatabase) DSN(t *testing.T) string {
	return fmt.Sprintf("postgres://postgres:postgres@127.0.0.1:%d/accounts_test?sslmode=disable", db.Port(t))
}

func (db *PostgresDatabase) Close(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	require.NoError(t, db.instance.Terminate(ctx))
}
