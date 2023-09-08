package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/alexzhirkov/url-shortener/internal/domain"
	"github.com/alexzhirkov/url-shortener/internal/repository/url_repository"
	"sync"

	_ "github.com/mattn/go-sqlite3" //init.go sqlite3 driver
)

type Storage struct {
	db *sql.DB
	mu sync.RWMutex
}

func New(storagePath string) (*Storage, error) {
	const fn = "adapters.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s. sql.Open: %w", fn, err)
	}

	stmt, err := db.Prepare(`
	create table if not exists http(
		id integer primary key,
		alias text not null unique,
		http text not null
	);
	create index if not exists idx_alias on http(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s, db.Prepare: %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s. stmt.Exec: %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (store *Storage) GetUrl(ctx context.Context, alias string) (*domain.ShortUrl, error) {
	const fn = "adapters.sqlite.GetUrl"

	var shortUrl url_repository.ShortUrl
	err := store.db.QueryRow(`select * from http where alias = ?;`, alias).Scan(&shortUrl.Id, &shortUrl.Alias, &shortUrl.Url)
	if err != nil {
		if errors.Is(errors.New("sql: no rows in result set"), err) {
			return &domain.ShortUrl{}, url_repository.ErrAliasNotFound
		}
		return &domain.ShortUrl{}, fmt.Errorf("%s. db.QueryRow: %w", fn, err)
	}
	storeUrl, err := domain.NewShortUrl(shortUrl.Id, shortUrl.Url, shortUrl.Alias)
	if err != nil {
		return &domain.ShortUrl{}, fmt.Errorf("%s. domain.NewShortUrl: %w", fn, err)
	}
	return storeUrl, nil

}

// Count returns the number of urls
func (store *Storage) Count() (int, error) {
	const fn = "adapters.sqlite.Count"

	rows, err := store.db.Query(`select count(*) cnt from http;`)
	if err != nil {
		return 0, fmt.Errorf("%s. db.Query: %w", fn, err)
	}
	defer rows.Close()
	var cnt int
	for rows.Next() {
		if err := rows.Scan(&cnt); err != nil {
			return 0, fmt.Errorf("%s. rows.Next: %w", fn, err)
		}
	}
	return cnt, nil
}

// CreateUrl creates or updates a port
func (store *Storage) CreateUrl(ctx context.Context, shortUrl *domain.ShortUrl) error {
	const fn = "adapters.sqlite.CreateUrl"

	store.mu.Lock()
	defer store.mu.Unlock()
	stmt, err := store.db.Prepare(`
	insert into http (id, alias, http) values(?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("%s, db.Prepare: %w", fn, err)
	}
	if _, err := stmt.Exec(shortUrl.GetId(), shortUrl.GetAlias(), shortUrl.GetUrl()); err != nil {
		return fmt.Errorf("%s. stmt.Exec: %w", fn, err)
	}

	return nil
}
