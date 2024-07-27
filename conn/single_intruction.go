package conn

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type SingleInstruction struct {
	db *sqlx.DB
}

func NewSingleInstruction(db *sqlx.DB) *SingleInstruction {
	return &SingleInstruction{
		db: db,
	}
}

func (s *SingleInstruction) Query(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return s.db.QueryxContext(ctx, query, args...)
}

func (s *SingleInstruction) QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return s.db.QueryRowxContext(ctx, query, args...)
}

func (s *SingleInstruction) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *SingleInstruction) Prepare(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return s.db.PreparexContext(ctx, query)
}

func (s *SingleInstruction) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.db.SelectContext(ctx, dest, query, args...)
}

func (s *SingleInstruction) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.db.GetContext(ctx, dest, query, args...)
}

func (s *SingleInstruction) Rebind(query string) string {
	return s.db.Rebind(query)
}

func (s *SingleInstruction) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return s.db.NamedExecContext(ctx, query, arg)
}
