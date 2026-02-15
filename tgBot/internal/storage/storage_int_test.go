package storage_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *redis.Client

func setupValkey(ctx context.Context) (testcontainers.Container, error) {
	redisConn, err := testcontainers.Run(
		ctx, "valkey/valkey:9.0.0",
		testcontainers.WithExposedPorts("6379/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("6379/tcp"),
			wait.ForLog("Ready to accept connections"),
		),
	)

	return redisConn, err
}

func TestValkeyConnection(t *testing.T) {
	const op = "TestRedisConnection"

	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	ctx := context.Background()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	valkey, err := setupValkey(ctx)

	defer func() {
		if err := testcontainers.TerminateContainer(valkey); err != nil {
			log.Error("Failed to terminate Redis container", "op", op, "err", err)
		}
	}()
	testcontainers.CleanupContainer(t, valkey)
	require.NoError(t, err)
}

func TestUserCache_SetTempUserState(t *testing.T) {
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	ctx := context.Background()
	valkeyContainer, _ := setupValkey(ctx)
	defer valkeyContainer.Terminate(ctx)

	valkey :=



}
