package duckdbmanager

import (
	"context"
	"database/sql"
)

type Config struct {
	Datasource string
}

type DuckdbManager struct {
	db *sql.DB
}

type ShutdownFunc func()

func NewDuckdbManager(config Config) (*DuckdbManager, ShutdownFunc, error) {
	db, err := sql.Open("duckdb", config.Datasource)
	if err != nil {
		return nil, func() {}, err
	} else {
		manager := &DuckdbManager{db: db}
		shutdownFunc := func() {
			if err := manager.Close(); err != nil {
				// Handle error if needed
			}
		}
		return manager, shutdownFunc, nil
	}

}

func (m *DuckdbManager) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return m.db.ExecContext(ctx, query, args...)
}

func (m *DuckdbManager) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, query, args...)
}

func (m *DuckdbManager) Close() error {
	return m.db.Close()
}
