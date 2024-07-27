package use_case

import (
	"context"

	"github.com/pudding-hack/backend/be-inventory/internal/model/item"
	"github.com/pudding-hack/backend/be-inventory/internal/model/unit"
	"github.com/pudding-hack/backend/lib"
)

type inventoryRepository interface {
	GetAll(ctx context.Context) (res []item.Item, err error)
	GetByID(ctx context.Context, id string) (item.Item, error)
	Create(ctx context.Context, item item.Item) error
}

type unitRepository interface {
	GetUnitById(ctx context.Context, id int) (unit.Unit, error)
	GetUnitByIds(ctx context.Context, ids []int) (res []unit.Unit, err error)
}

type service struct {
	cfg      *lib.Config
	repo     inventoryRepository
	unitRepo unitRepository
}

func NewService(cfg *lib.Config, repo inventoryRepository, unitRepo unitRepository) *service {
	return &service{
		cfg:      cfg,
		repo:     repo,
		unitRepo: unitRepo,
	}
}

func (s *service) GetAll(ctx context.Context) (res []Item, err error) {
	inventories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	unitIds := []int{}
	for _, inventory := range inventories {
		unitIds = append(unitIds, inventory.UnitId)
	}

	units, err := s.unitRepo.GetUnitByIds(ctx, unitIds)
	if err != nil {
		return
	}
	unitMap := map[int]unit.Unit{}
	for _, unit := range units {
		unitMap[unit.ID] = unit
	}

	for _, inventory := range inventories {
		var i Item
		i.FromEntity(inventory, unitMap[inventory.UnitId].UnitName)
		res = append(res, i)
	}

	return res, nil
}

func (s *service) GetByID(ctx context.Context, id string) (Item, error) {
	inventory, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Item{}, err
	}

	unit, err := s.unitRepo.GetUnitById(ctx, inventory.UnitId)
	if err != nil {
		return Item{}, err
	}

	var i Item
	i.FromEntity(inventory, unit.UnitName)

	return i, nil
}

func (s *service) Create(ctx context.Context, item item.Item) error {
	err := s.repo.Create(ctx, item)
	if err != nil {
		return err
	}

	return nil
}
