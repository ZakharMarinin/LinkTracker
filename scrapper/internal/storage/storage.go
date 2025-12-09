package storage

import (
	"context"
	"embed"
	"fmt"
	"scrapper/internal/domain"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type PostgresStorage struct {
	DB *pgxpool.Pool
}

func New(ctx context.Context, storagePath string) (*PostgresStorage, error) {
	const op = "storage.postgresql.SQL.NEW"

	if err := runMigrations(storagePath); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	poolConfig, err := pgxpool.ParseConfig(storagePath)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{DB: db}, nil
}

func runMigrations(dbURL string) error {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return err
	}

	db := stdlib.OpenDB(*config.ConnConfig)
	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) Close() error {
	s.DB.Close()
	return nil
}

func (s *PostgresStorage) CreateChat(ctx context.Context, chatID int64) error {
	const op = "storage.postgres.CreateChat"

	create, args, err := sq.
		Insert("users").
		Columns("chat_id").
		Values(chatID).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.DB.Exec(ctx, create, args...)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil
}

func (s *PostgresStorage) DeleteChat(ctx context.Context, chatID int64) error {
	const op = "storage.postgres.DeleteChat"

	deleteQ, args, err := sq.
		Delete("users").
		Where(sq.Eq{"chat_id": chatID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.DB.Exec(ctx, deleteQ, args...)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil
}

func (s *PostgresStorage) GetLinksToCheck(ctx context.Context, limit, offset uint64) ([]domain.Link, error) {
	const op = "storage.postgres.GetLinksToCheck"

	query, args, err := sq.
		Select("l.chat_id", "l.link_id", "ul.url", "ul.alias", "ul.last_update").
		From("links ul").
		Join("user_links l ON ul.link_id = l.link_id").
		Limit(limit).
		Offset(offset).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var links []domain.Link
	for rows.Next() {
		var link domain.Link
		if err := rows.Scan(&link.ChatID, &link.ID, &link.URL, &link.Alias, &link.LastUpdated); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return links, nil
}

func (s *PostgresStorage) GetLinks(ctx context.Context, limit, offset uint64) ([]domain.Link, error) {
	const op = "storage.postgres.GetLinks"

	query, _, err := sq.
		Select("*").
		From("links").
		Limit(limit).
		Offset(offset).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var links []domain.Link

	for rows.Next() {
		var link domain.Link
		if err := rows.Scan(&link.ID, &link.URL, &link.Alias, &link.LastUpdated); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return links, nil
}

func (s *PostgresStorage) GetUserLinksByTag(ctx context.Context, chatID int64, tags string) ([]*domain.Link, error) {
	const op = "storage.postgres.GetUserLinksByTag"

	query, args, err := sq.
		Select("l.link_id", "l.url", "ul.alias", "ul.description", "ul.tags").
		From("user_links ul").
		Join("links l ON ul.link_id = l.link_id").
		Where(sq.And{
			sq.Eq{"ul.chat_id": chatID},
			sq.Expr("ul.tags LIKE ?", "%"+tags+"%"),
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()
	var links []*domain.Link
	for rows.Next() {
		var link domain.Link
		err := rows.Scan(&link.ID, &link.URL, &link.Alias, &link.Desc, &link.Tags)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		links = append(links, &link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return links, nil
}

func (s *PostgresStorage) GetLinksByChatID(ctx context.Context, chatID int64) ([]domain.Link, error) {
	const op = "storage.postgres.GetLinksByChatID"

	query, args, err := sq.
		Select("l.link_id", "l.url", "ul.alias", "ul.description", "ul.tags").
		From("user_links ul").Join("links l ON ul.link_id = l.link_id").
		Where(sq.Eq{"ul.chat_id": chatID}).
		OrderBy("l.link_id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var links []domain.Link

	for rows.Next() {
		var link domain.Link
		if err := rows.Scan(&link.ID, &link.URL, &link.Alias, &link.Desc, &link.Tags); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return links, nil
}

func (s *PostgresStorage) GetLinkIDByURL(ctx context.Context, url string) (*domain.Link, error) {
	const op = "storage.postgres.GetLinkByURL"

	query, args, err := sq.
		Select("link_id").
		From("links").
		Where(sq.Eq{"url": url}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows := s.DB.QueryRow(ctx, query, args...)

	var link domain.Link

	if err := rows.Scan(&link.ID); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &link, nil
}

func (s *PostgresStorage) GetLinkByAlias(ctx context.Context, chatID int64, alias string) (*domain.Link, error) {
	const op = "storage.postgres.GetLinkByAlias"

	query, args, err := sq.
		Select("link_id").
		From("user_links").
		Where(sq.Eq{"chat_id": chatID, "alias": alias}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows := s.DB.QueryRow(ctx, query, args...)

	var link domain.Link

	if err := rows.Scan(&link.ID); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &link, nil
}

func (s *PostgresStorage) GetLinkByURL(ctx context.Context, url string) (*domain.Link, error) {
	const op = "storage.postgres.GetLinkByURL"

	query, args, err := sq.
		Select("link_id", "url", "alias", "last_update").
		From("links").
		Where(sq.Eq{"url": url}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var link domain.Link
	err = s.DB.QueryRow(ctx, query, args...).Scan(&link.ID, &link.URL, &link.Alias, &link.LastUpdated)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &link, nil
}

func (s *PostgresStorage) IsLinkExists(ctx context.Context, url string) (bool, error) {
	const op = "storage.postgres.IsLinkExists"

	query, args, err := sq.
		Select("1").
		Prefix("SELECT EXISTS (").
		From("links").
		Where(sq.Eq{"url": url}).
		PlaceholderFormat(sq.Dollar).
		Suffix(")").
		ToSql()

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	var exists bool
	err = s.DB.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

func (s *PostgresStorage) IsUserLinkExists(ctx context.Context, alias string, chatID int64) (bool, error) {
	const op = "storage.postgres.IsLinkExists"

	query, args, err := sq.
		Select("1").
		Prefix("SELECT EXISTS (").
		From("user_links").
		Where(sq.Eq{"chat_id": chatID, "alias": alias}).
		PlaceholderFormat(sq.Dollar).
		Suffix(")").
		ToSql()

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	var exists bool
	err = s.DB.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

func (s *PostgresStorage) AddLink(ctx context.Context, link *domain.Link) error {
	const op = "storage.postgres.AddLink"

	curTime := time.Now()

	query, args, err := sq.
		Insert("links").
		Columns("url", "alias", "last_update").
		Values(link.URL, link.Alias, curTime).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *PostgresStorage) AddUserLink(ctx context.Context, chatID int64, link *domain.Link) error {
	const op = "storage.postgres.AddUserLink"

	intChatID := int(chatID)

	query, args, err := sq.
		Insert("user_links").
		Columns("chat_id", "link_id", "alias", "description", "tags").
		Values(intChatID, link.ID, link.Alias, link.Desc, link.Tags).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *PostgresStorage) UpdateLink(ctx context.Context, link *domain.Link) (*domain.Link, error) {
	const op = "storage.postgres.UpdateLink"

	var updatedLink domain.Link

	curTime := time.Now()

	query, args, err := sq.
		Update("links").
		Set("last_update", curTime).
		Where(sq.Eq{"url": link.URL}).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING link_id, url, last_update").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = s.DB.QueryRow(ctx, query, args...).Scan(&updatedLink.ID, &updatedLink.URL, &updatedLink.LastUpdated)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &updatedLink, nil
}

func (s *PostgresStorage) DeleteUserLink(ctx context.Context, chatID int64, alias string) error {
	const op = "storage.postgres.DeleteLink"

	query, args, err := sq.
		Delete("user_links").
		Where(sq.Eq{"chat_id": chatID, "alias": alias}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *PostgresStorage) DeleteLink(ctx context.Context, link *domain.Link) error {
	const op = "storage.postgres.DeleteLink"

	query, args, err := sq.
		Delete("links").
		Where(sq.Eq{"url": link.URL}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
