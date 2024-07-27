package conn

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Transaction interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type MultiInstruction struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewMultiInstruction(db *sqlx.DB) *MultiInstruction {
	return &MultiInstruction{
		db: db,
		tx: nil,
	}
}

func (t *MultiInstruction) Begin(ctx context.Context) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	t.tx = tx
	return nil
}

func (t *MultiInstruction) Commit(ctx context.Context) error {
	err := t.tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (t *MultiInstruction) Rollback(ctx context.Context) error {
	err := t.tx.Rollback()
	if err != nil {
		return err
	}

	return nil
}

func (t *MultiInstruction) Query(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return t.tx.QueryxContext(ctx, query, args...)
}

func (t *MultiInstruction) QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return t.tx.QueryRowxContext(ctx, query, args...)
}

func (t *MultiInstruction) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *MultiInstruction) Prepare(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return t.tx.PreparexContext(ctx, query)
}

func (t *MultiInstruction) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return t.tx.SelectContext(ctx, dest, query, args...)
}

func (t *MultiInstruction) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return t.tx.GetContext(ctx, dest, query, args...)
}

func (t *MultiInstruction) CommitAndClose(ctx context.Context) error {
	err := t.tx.Commit()
	if err != nil {
		return err
	}

	t.tx = nil
	return nil
}

func (t *MultiInstruction) RollbackAndClose(ctx context.Context) error {
	err := t.tx.Rollback()
	if err != nil {
		return err
	}

	t.tx = nil
	return nil
}

func (t *MultiInstruction) Rebind(query string) string {
	return t.tx.Rebind(query)
}

func (t *MultiInstruction) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return t.tx.NamedExecContext(ctx, query, arg)
}
