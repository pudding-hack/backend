package use_case

import (
	"context"

	"github.com/pudding-hack/backend/be-inventory/internal/model/history"
	"github.com/pudding-hack/backend/be-inventory/internal/model/item"
	"github.com/pudding-hack/backend/be-inventory/internal/model/unit"
	"github.com/pudding-hack/backend/lib"
)

type inventoryRepository interface {
	GetAll(ctx context.Context) (res []item.Item, err error)
	GetByID(ctx context.Context, id int) (item.Item, error)
	Create(ctx context.Context, item item.Item) error
}

type unitRepository interface {
	GetUnitById(ctx context.Context, id int) (unit.Unit, error)
	GetUnitByIds(ctx context.Context, ids []int) (res []unit.Unit, err error)
}

type historyRepository interface {
	GetByID(ctx context.Context, id string, request lib.PaginationRequest) (res []history.HistoryItem, meta lib.Pagination, err error)
	CreateHistory(ctx context.Context, item history.HistoryItem) error
	GetHistoryTypeByIds(ctx context.Context, ids []int) (res []history.HistoryType, err error)
}

type service struct {
	cfg      *lib.Config
	repo     inventoryRepository
	unitRepo unitRepository
	histRepo historyRepository
}

func NewService(cfg *lib.Config, repo inventoryRepository, unitRepo unitRepository, histRepo historyRepository) *service {
	return &service{
		cfg:      cfg,
		repo:     repo,
		unitRepo: unitRepo,
		histRepo: histRepo,
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

func (s *service) GetByID(ctx context.Context, id int) (Item, error) {
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

func (s *service) GetItemHistoryPaginate(ctx context.Context, id string, request lib.PaginationRequest) (response GetHistoryResponse, err error) {

	histories, meta, err := s.histRepo.GetByID(ctx, id, request)
	if err != nil {
		return
	}

	historiesTypeIds := []int{}
	for _, history := range histories {
		historiesTypeIds = append(historiesTypeIds, history.TypeId)
	}

	historiesType, err := s.histRepo.GetHistoryTypeByIds(ctx, historiesTypeIds)
	if err != nil {
		return
	}

	historiesTypeMap := map[int]history.HistoryType{}
	for _, historyType := range historiesType {
		historiesTypeMap[historyType.ID] = historyType
	}

	for _, history := range histories {
		var h HistoryItem
		h.FromEntity(history, historiesTypeMap[history.TypeId].TypeName)
		response.Data = append(response.Data, h)
	}

	response.Meta = meta

	return response, nil
}

func (s *service) InboundItem(ctx context.Context, id int, qty int) (err error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return
	}

	historyItem := history.HistoryItem{
		ItemId:   item.ID,
		Quantity: qty,
		TypeId:   1,
	}

	err = s.histRepo.CreateHistory(ctx, historyItem)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) OutboundItem(ctx context.Context, id int, qty int) (err error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return
	}

	historyItem := history.HistoryItem{
		ItemId:   item.ID,
		Quantity: qty,
		TypeId:   2,
	}

	err = s.histRepo.CreateHistory(ctx, historyItem)
	if err != nil {
		return err
	}

	return nil
}
