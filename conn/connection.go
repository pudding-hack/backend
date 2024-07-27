package conn

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Connection interface {
	Query(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(ctx context.Context, query string) (*sqlx.Stmt, error)
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Rebind(query string) string
	NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}
