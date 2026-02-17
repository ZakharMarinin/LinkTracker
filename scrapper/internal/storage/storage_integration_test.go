package storage_test

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"scrapper/internal/domain"
	"scrapper/internal/storage"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresTestSuite struct {
	suite.Suite
	repo    *storage.PostgresStorage
	cleanup func()
}

func SetupPostgres(ctx context.Context) (*storage.PostgresStorage, func(), error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "test",
		},
		WaitingFor: wait.ForSQL("5432/tcp", "pgx", func(host string, port nat.Port) string {
			return fmt.Sprintf("postgres://testuser:password@%s:%s/test?sslmode=disable", host, port.Port())
		}).WithStartupTimeout(time.Minute),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start container: %w", err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("postgres://testuser:password@%s:%s/test?sslmode=disable", host, port.Port())

	migrationDB, err := sql.Open("pgx", dsn)
	if err != nil {
		err = container.Terminate(ctx)
		if err != nil {
			return nil, nil, err
		}

		return nil, nil, fmt.Errorf("failed to open db for migrations: %w", err)
	}

	_, filename, _, _ := runtime.Caller(0)
	migrationPath := filepath.Join(filepath.Dir(filename), "../../../migrations")

	if err = goose.SetDialect("postgres"); err != nil {
		err = migrationDB.Close()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to close migrations db: %w", err)
		}

		err = container.Terminate(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to terminate container: %w", err)
		}

		return nil, nil, err
	}

	if err = goose.Up(migrationDB, migrationPath); err != nil {
		err = migrationDB.Close()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to close migrations db: %w", err)
		}

		err = container.Terminate(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to terminate dicontainer: %w", err)
		}

		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	err = migrationDB.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to close migrations db: %w", err)
	}

	pool, err := storage.New(ctx, dsn)
	if err != nil {
		return nil, nil, err
	}

	if err = pool.DB.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("unable to ping database: %s", err)
	}

	cleanup := func() {
		pool.DB.Close()

		err = container.Terminate(ctx)
		if err != nil {
			return
		}
	}

	return pool, cleanup, nil
}

func (s *PostgresTestSuite) SetupSuite() {
	db, cleanup, err := SetupPostgres(context.Background())
	s.Require().NoError(err)
	s.repo = db
	s.cleanup = cleanup
}

func (s *PostgresTestSuite) TearDownSuite() {
	s.cleanup()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

func (s *PostgresTestSuite) Test_CreateChatAndDeleteChat() {
	chatID := int64(1111)

	err := s.repo.CreateChat(context.Background(), chatID)
	s.NoError(err)

	exists := s.checkChatExists(chatID)
	s.True(exists)

	err = s.repo.DeleteChat(context.Background(), chatID)
	s.NoError(err)

	exists = s.checkChatExists(chatID)
	s.False(exists)
}

func (s *PostgresTestSuite) Test_CreateChatDuplicate() {
	chatID := int64(1111)

	err := s.repo.CreateChat(context.Background(), chatID)
	s.NoError(err)

	exists := s.checkChatExists(chatID)
	s.True(exists)

	err = s.repo.CreateChat(context.Background(), chatID)

	s.Require().Error(err, "duplicate chat")
}

func (s *PostgresTestSuite) checkChatExists(chatID int64) bool {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE chat_id = $1)"
	err := s.repo.DB.QueryRow(context.Background(), query, chatID).Scan(&exists)
	s.NoError(err)

	return exists
}

func (s *PostgresTestSuite) Test_CreateLinkAndDeleteLink() {
	link := &domain.Link{
		URL:         "https://google.com",
		Alias:       "google",
		LastUpdated: time.Now(),
	}

	err := s.repo.AddLink(context.Background(), link)
	s.Require().NoError(err)

	exists, err := s.repo.IsLinkExists(context.Background(), link.URL)
	s.Require().NoError(err)
	s.True(exists)

	err = s.repo.DeleteLink(context.Background(), link)
	s.Require().NoError(err)

	exists, err = s.repo.IsLinkExists(context.Background(), link.URL)
	s.Require().NoError(err)
	s.False(exists)
}

func (s *PostgresTestSuite) Test_CreateUserLinkAndDeleteUserLink() {
	chatID := int64(1111)
	link := &domain.Link{
		URL:         "https://google.com",
		Alias:       "google",
		Desc:        "google",
		Tags:        "",
		LastUpdated: time.Now(),
	}

	err := s.repo.AddLink(context.Background(), link)
	s.Require().NoError(err, "failed to add link")

	exists, err := s.repo.IsLinkExists(context.Background(), link.URL)
	s.Require().NoError(err, "failed to check if link exists")
	s.Require().True(exists)

	link, err = s.repo.GetLinkByURL(context.Background(), link.URL)
	s.Require().NoError(err, "failed to get link by url")

	err = s.repo.AddUserLink(context.Background(), chatID, link)
	s.Require().NoError(err, "failed to add user link")

	exists, err = s.repo.IsUserLinkExists(context.Background(), link.Alias, chatID)
	s.Require().NoError(err, "failed to check user link")
	s.Require().True(exists)

	err = s.repo.DeleteUserLink(context.Background(), chatID, link.Alias)
	s.Require().NoError(err, "failed to delete user link")

	exists, err = s.repo.IsUserLinkExists(context.Background(), link.Alias, chatID)
	s.Require().NoError(err, "failed to check user link")
	s.Require().False(exists)
}

func (s *PostgresTestSuite) Test_GetLinksToCheck() {
	limit := uint64(100)
	offset := uint64(0)

	for range 10 {
		url := fmt.Sprintf("https://%s.com", uuid.New().String())

		err := s.repo.AddLink(context.Background(), &domain.Link{
			URL:         url,
			Alias:       "google",
			Desc:        "google",
			Tags:        "",
			LastUpdated: time.Now(),
		})
		s.Require().NoError(err, "failed to add link")

		exists, err := s.repo.IsLinkExists(context.Background(), url)
		s.Require().NoError(err, "failed to check if link exists")
		s.Require().True(exists)

		link, err := s.repo.GetLinkByURL(context.Background(), url)
		s.Require().NoError(err, "failed to get link by url")

		chatID := int64(1111)

		err = s.repo.AddUserLink(context.Background(), chatID, link)
		s.Require().NoError(err, "failed to add user link")

		exists, err = s.repo.IsUserLinkExists(context.Background(), link.Alias, chatID)
		s.Require().NoError(err, "failed to check user link")
		s.Require().True(exists)
	}

	links, err := s.repo.GetLinksToCheck(context.Background(), limit, offset)
	s.Require().NoError(err, "failed to get links to check")
	s.Require().Len(links, 10)
}
