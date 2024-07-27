package unit

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pudding-hack/backend/conn"
	"github.com/pudding-hack/backend/lib"
)

type repository struct {
	cfg *lib.Config
	db  conn.Connection
}

func New(cfg *lib.Config, db conn.Connection) *repository {
	return &repository{
		cfg: cfg,
		db:  db,
	}
}

func (r *repository) GetUnitById(ctx context.Context, id int) (Unit, error) {
	var unit Unit
	err := r.db.Get(ctx, &unit, "SELECT * FROM units WHERE id = $1 AND deleted_at is NULL", id)
	if err != nil {
		return Unit{}, err
	}

	return unit, nil
}

func (r *repository) GetUnitByIds(ctx context.Context, ids []int) (res []Unit, err error) {
	var unit []Unit
	query, args, err := sqlx.In("SELECT * FROM units WHERE id IN (?) AND deleted_at is NULL", ids)
	if err != nil {
		return
	}

	err = r.db.Select(ctx, &unit, r.db.Rebind(query), args...)
	if err != nil {
		return
	}

	return unit, nil
}
