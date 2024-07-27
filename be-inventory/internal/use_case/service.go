package use_case

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/pudding-hack/backend/be-inventory/internal/model/history"
	"github.com/pudding-hack/backend/be-inventory/internal/model/item"
	"github.com/pudding-hack/backend/be-inventory/internal/model/unit"
	"github.com/pudding-hack/backend/conn"
	"github.com/pudding-hack/backend/lib"
)

type inventoryRepository interface {
	GetAll(ctx context.Context) (res []item.Item, err error)
	GetByID(ctx context.Context, id int) (item.Item, error)
	Create(ctx context.Context, item item.Item) error
	UpdateQuantity(ctx context.Context, id, qty int) error
	GetByName(ctx context.Context, name string) (item.Item, error)
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

type awsService interface {
	DetectLabels(ctx context.Context, params *rekognition.DetectLabelsInput, optFns ...func(*rekognition.Options)) (*rekognition.DetectLabelsOutput, error)
}

type service struct {
	cfg        *lib.Config
	repo       inventoryRepository
	unitRepo   unitRepository
	histRepo   historyRepository
	conn       *conn.SQLServerConnectionManager
	db         *conn.SingleInstruction
	tx         *conn.MultiInstruction
	awsService awsService
}

func NewService(sql *conn.SQLServerConnectionManager, cfg *lib.Config, awsService awsService) *service {
	return &service{
		cfg:        cfg,
		db:         sql.GetQuery(),
		conn:       sql,
		repo:       item.New(cfg, sql.GetQuery()),
		unitRepo:   unit.New(cfg, sql.GetQuery()),
		histRepo:   history.New(cfg, sql.GetQuery()),
		awsService: awsService,
	}
}

func (s *service) WithTransaction() {
	s.tx = s.conn.GetTransaction()
	s.repo = item.New(s.cfg, s.tx)
	s.unitRepo = unit.New(s.cfg, s.tx)
	s.histRepo = history.New(s.cfg, s.tx)
}

func (s *service) WithoutTransaction() {
	s.repo = item.New(s.cfg, s.db)
	s.unitRepo = unit.New(s.cfg, s.db)
	s.histRepo = history.New(s.cfg, s.db)
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
		i.FromEntity(inventory)
		i.Unit = unitMap[inventory.UnitId].UnitName
		res = append(res, i)
	}

	return res, nil
}

func (s *service) GetByID(ctx context.Context, id int) (Item, error) {
	inventory, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, lib.ErrSqlErrorNotFound) {
			return Item{}, lib.NewErrNotFound("Item not found")
		}

		return Item{}, err
	}

	unit, err := s.unitRepo.GetUnitById(ctx, inventory.UnitId)
	if err != nil {
		if errors.Is(err, lib.ErrSqlErrorNotFound) {
			return Item{}, lib.NewErrNotFound("Unit not found")
		}

		return Item{}, err
	}

	var i Item
	i.FromEntity(inventory)
	i.Unit = unit.UnitName

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

	if len(histories) == 0 {
		return response, lib.NewErrNotFound("History not found")
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

func (s *service) InboundItem(ctx context.Context, name string, qty int) (err error) {
	s.WithTransaction()

	defer func() {
		if p := recover(); p != nil {
			s.tx.Rollback(ctx)
			err = fmt.Errorf("panic: %v", p)
		} else if err != nil {
			s.tx.Rollback(ctx)
		} else {
			log.Println("Commit")
			s.tx.Commit(ctx)
		}
		s.WithoutTransaction()
	}()

	s.tx.Begin(ctx)

	user := lib.GetUserContext(ctx)

	item, err := s.repo.GetByName(ctx, name)
	if err != nil {
		if errors.Is(err, lib.ErrSqlErrorNotFound) {
			return lib.NewErrNotFound("Item not found")
		}
	}

	historyItem := history.HistoryItem{
		ItemId:         item.ID,
		Quantity:       qty,
		TypeId:         1,
		QuantityBefore: item.Quantity,
		QuantityAfter:  item.Quantity + qty,
		CreatedBy:      user.ID,
		UpdatedBy:      user.ID,
	}

	err = s.histRepo.CreateHistory(ctx, historyItem)
	if err != nil {
		if errors.Is(err, lib.ErrSqlErrorNotFound) {
			return lib.NewErrNotFound("Item not found")
		}
		return err
	}

	err = s.repo.UpdateQuantity(ctx, item.ID, qty)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) OutboundItem(ctx context.Context, name string, qty int) (err error) {
	s.WithTransaction()

	defer func() {
		if p := recover(); p != nil {
			s.tx.Rollback(ctx)
			err = fmt.Errorf("panic: %v", p)
		} else if err != nil {
			s.tx.Rollback(ctx)
		} else {
			log.Println("Commit")
			s.tx.Commit(ctx)
		}
		s.WithoutTransaction()
	}()

	s.tx.Begin(ctx)

	user := lib.GetUserContext(ctx)

	item, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return
	}

	historyItem := history.HistoryItem{
		ItemId:         item.ID,
		Quantity:       qty,
		TypeId:         2,
		QuantityBefore: item.Quantity,
		QuantityAfter:  item.Quantity - qty,
		CreatedBy:      user.ID,
		UpdatedBy:      user.ID,
	}

	err = s.histRepo.CreateHistory(ctx, historyItem)
	if err != nil {
		return err
	}

	err = s.repo.UpdateQuantity(ctx, item.ID, -qty)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DetectLabels(ctx context.Context, imageBase64 string) (res Item, err error) {
	image, err := lib.ConvertBase64ToImage(imageBase64)
	if err != nil {
		return
	}

	params := &rekognition.DetectLabelsInput{
		Image: &types.Image{
			Bytes: image,
		},
	}

	resp, err := s.awsService.DetectLabels(ctx, params)
	if err != nil {
		return
	}

	if len(resp.Labels) == 0 {
		return res, lib.NewErrNotFound("Label not found")
	}

	var label string

	for _, l := range resp.Labels {
		if *l.Confidence > 90 {
			label = *l.Name
			break
		}
	}

	if label == "" {
		return res, lib.NewErrNotFound("Label not found")
	}

	item, err := s.repo.GetByName(ctx, label)
	if err != nil {
		if errors.Is(err, lib.ErrSqlErrorNotFound) {
			return Item{}, lib.NewErrNotFound(label + " not found")
		}
		return
	}

	unit, err := s.unitRepo.GetUnitById(ctx, item.UnitId)
	if err != nil {
		return
	}

	res.FromEntity(item)
	res.Unit = unit.UnitName

	return res, nil
}
