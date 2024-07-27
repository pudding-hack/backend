package history

import "time"

type HistoryItem struct {
	ID             int        `json:"id" db:"id"`
	ItemId         int        `json:"item_id" db:"item_id"`
	Quantity       int        `json:"qty" db:"qty"`
	QuantityBefore int        `json:"qty_before" db:"quantity_before"`
	QuantityAfter  int        `json:"qty_after" db:"quantity_after"`
	TypeId         int        `json:"type_id" db:"type_id"`
	Note           *string    `json:"note" db:"note"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy      string     `json:"created_by" db:"created_by"`
	UpdatedBy      string     `json:"updated_by" db:"updated_by"`
	DeletedBy      *string    `json:"deleted_by" db:"deleted_by"`
}

type HistoryType struct {
	ID        int        `json:"id" db:"id"`
	TypeName  string     `json:"name" db:"name"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy string     `json:"created_by" db:"created_by"`
	UpdatedBy string     `json:"updated_by" db:"updated_by"`
	DeletedBy *string    `json:"deleted_by" db:"deleted_by"`
}
