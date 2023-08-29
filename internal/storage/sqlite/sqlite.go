package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" //init sqlite3 driver
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s. sql.Open: %w", fn, err)
	}

	stmt, err := db.Prepare(`
	create table if not exists url(
		id integer primary key,
		alias text not null unique,
		url text not null
	);
	create index if not exists idx_alias on url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s, db.Prepare: %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s. stmt.Exec: %w", fn, err)
	}

	return &Storage{db}, nil
}
