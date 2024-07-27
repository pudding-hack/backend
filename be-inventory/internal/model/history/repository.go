package history

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

func (r *repository) GetByID(ctx context.Context, id string, request lib.PaginationRequest) (res []HistoryItem, meta lib.Pagination, err error) {
	var histories []HistoryItem

	query := "SELECT * FROM history_items WHERE deleted_at IS NULL AND item_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"

	err = r.db.Select(ctx, &histories, query, id, request.Limit(), request.Offset())

	if err != nil {
		return
	}

	query = `SELECT COUNT(*) FROM history_items WHERE item_id = $1`
	var total int
	err = r.db.Get(ctx, &total, query, id)
	if err != nil {
		return
	}

	meta = lib.NewPagination(request.Page, request.PageSize, total)

	return histories, meta, nil
}

func (r *repository) GetHistoryTypeByIds(ctx context.Context, ids []int) (res []HistoryType, err error) {
	var types []HistoryType
	query, args, err := sqlx.In("SELECT * FROM history_types WHERE id IN (?)", ids)
	if err != nil {
		return
	}

	err = r.db.Select(ctx, &types, r.db.Rebind(query), args...)
	if err != nil {
		return
	}

	return types, nil
}

func (r *repository) CreateHistory(ctx context.Context, item HistoryItem) error {
	query := `INSERT INTO history_items (item_id, type_id, qty, created_by, updated_by) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, query, item.ItemId, item, item.Quantity, item.CreatedBy, item.UpdatedBy)
	if err != nil {
		return err
	}

	return nil
}
