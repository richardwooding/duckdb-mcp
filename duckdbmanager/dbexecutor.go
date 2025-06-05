package duckdbmanager

import (
	"context"
	"database/sql"
)

type DbExecutor interface {
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
}
