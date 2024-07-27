package item

import (
	"time"
)

type Item struct {
	ID        int        `json:"id" db:"id"`
	ItemCode  string     `json:"item_code" db:"item_code"`
	ItemName  string     `json:"item_name" db:"item_name"`
	Quantity  int        `json:"qty" db:"qty"`
	UnitId    int        `json:"unit" db:"unit"`
	Price     float64    `json:"price" db:"price"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy string     `json:"created_by" db:"created_by"`
	UpdatedBy string     `json:"updated_by" db:"updated_by"`
	DeletedBy *string    `json:"deleted_by" db:"deleted_by"`
}
