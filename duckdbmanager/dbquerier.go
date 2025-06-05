package duckdbmanager

import (
	"context"
	"database/sql"
)

type DbQuerier interface {
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}
