package use_case

import (
	"context"

	"github.com/pudding-hack/backend/be-inventory/internal/model/item"
	"github.com/pudding-hack/backend/lib"
)

type inventoryRepository interface {
	GetAll(ctx context.Context) (res []item.Item, err error)
	GetByID(ctx context.Context, id string) (item.Item, error)
	Create(ctx context.Context, item item.Item) error
}

type service struct {
	cfg  *lib.Config
	repo inventoryRepository
}

func NewService(cfg *lib.Config, repo inventoryRepository) *service {
	return &service{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *service) GetAll(ctx context.Context) (res []Item, err error) {
	inventories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, inventory := range inventories {
		var i Item
		i.FromEntity(inventory)
		res = append(res, i)
	}

	return res, nil
}

func (s *service) GetByID(ctx context.Context, id string) (Item, error) {
	inventory, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Item{}, err
	}

	var i Item
	i.FromEntity(inventory)

	return i, nil
}

func (s *service) Create(ctx context.Context, item item.Item) error {
	err := s.repo.Create(ctx, item)
	if err != nil {
		return err
	}

	return nil
}
