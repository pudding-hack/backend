package use_case

import (
	"github.com/pudding-hack/backend/be-inventory/internal/model/history"
	"github.com/pudding-hack/backend/be-inventory/internal/model/item"
	"github.com/pudding-hack/backend/lib"
)

type Item struct {
	ID       int     `json:"id"`
	ItemCode string  `json:"item_code"`
	ItemName string  `json:"item_name"`
	Quantity int     `json:"qty"`
	UnitId   int     `json:"unit_id"`
	Unit     string  `json:"unit"`
	Price    float64 `json:"price"`
}

func (i *Item) FromEntity(entity item.Item) {
	i.ID = entity.ID
	i.ItemCode = entity.ItemCode
	i.ItemName = entity.ItemName
	i.Quantity = entity.Quantity
	i.UnitId = entity.UnitId
	i.Price = entity.Price
}

type GetHistoryResponse struct {
	Data []HistoryItem  `json:"data"`
	Meta lib.Pagination `json:"meta"`
}

type HistoryItem struct {
	ID             int    `json:"id"`
	ItemId         int    `json:"item_id"`
	Quantity       int    `json:"qty"`
	QuantityBefore int    `json:"qty_before"`
	QuantityAfter  int    `json:"qty_after"`
	TypeId         int    `json:"type_id"`
	Type           string `json:"type"`
	CreatedAt      string `json:"created_at"`
}

func (h *HistoryItem) FromEntity(entity history.HistoryItem, typeName string) {
	h.ID = entity.ID
	h.ItemId = entity.ItemId
	h.Quantity = entity.Quantity
	h.TypeId = entity.TypeId
	h.Type = typeName
	h.QuantityBefore = entity.QuantityBefore
	h.QuantityAfter = entity.QuantityAfter
	h.CreatedAt = entity.CreatedAt.Format("02 January 2006 15:04:05")
}
