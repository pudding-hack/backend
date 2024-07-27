package item

import (
	"context"

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

func (r *repository) GetByName(ctx context.Context, name string) (Item, error) {
	var item Item
	err := r.db.Get(ctx, &item, "SELECT * FROM items WHERE item_name ILIKE $1 AND deleted_at is NULL", "%"+name+"%")
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func (r *repository) GetAll(ctx context.Context) (res []Item, err error) {
	var item []Item
	err = r.db.Select(ctx, &item, "SELECT * FROM items WHERE deleted_at is NULL")
	if err != nil {
		return []Item{}, err
	}

	return item, nil
}

func (r *repository) GetByID(ctx context.Context, id int) (Item, error) {
	var item Item
	err := r.db.Get(ctx, &item, "SELECT * FROM items WHERE id = $1 AND deleted_at is NULL", id)
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func (r *repository) Create(ctx context.Context, item Item) error {
	query := "INSERT INTO items (item_code, item_name, qty, unit, price, keywords, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	err := r.db.Get(ctx, &item.ID, query, item.ItemCode, item.ItemName, item.Quantity, item.UnitId, item.Price, item.Keywords, item.CreatedBy, item.UpdatedBy)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateQuantity(ctx context.Context, id, qty int) error {
	user := lib.GetUserContext(ctx)

	query := "UPDATE items SET qty = qty + $1, updated_at = NOW(), updated_by = $2 WHERE id = $3"
	_, err := r.db.Exec(ctx, query, qty, user.ID, id)
	if err != nil {
		return err
	}

	return nil
}
