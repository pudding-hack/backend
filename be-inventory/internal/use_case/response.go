package use_case

import (
	"github.com/pudding-hack/backend/be-inventory/internal/model/item"
)

type Item struct {
	ID       string  `json:"id"`
	ItemCode string  `json:"item_code"`
	ItemName string  `json:"item_name"`
	Quantity int     `json:"qty"`
	UnitId   int     `json:"unit"`
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
