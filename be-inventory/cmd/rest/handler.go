package rest

import (
	"context"
	"net/http"

	"github.com/pudding-hack/backend/be-inventory/internal/model/item"
	"github.com/pudding-hack/backend/be-inventory/internal/use_case"
	"github.com/pudding-hack/backend/lib"
)

type service interface {
	GetAll(ctx context.Context) (res []use_case.Item, err error)
	GetByID(ctx context.Context, id int) (use_case.Item, error)
	Create(ctx context.Context, item item.Item) error
	GetItemHistoryPaginate(ctx context.Context, id string, request lib.PaginationRequest) (response use_case.GetHistoryResponse, err error)
	InboundItem(ctx context.Context, name string, qty int) (err error)
	OutboundItem(ctx context.Context, name string, qty int) (err error)
	DetectLabels(ctx context.Context, imageBase64 string) (res use_case.Item, err error)
}

type Handler struct {
	service service
}

func NewHandler(service service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	inventories, err := h.service.GetAll(ctx)
	if err != nil {
		lib.WriteResponse(w, err, nil)
		return
	}

	lib.WriteResponse(w, nil, inventories)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetQueryInt(r, "id", 0)

	inventory, err := h.service.GetByID(ctx, id)
	if err != nil {
		lib.WriteResponse(w, err, nil)
		return
	}

	lib.WriteResponse(w, nil, inventory)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var item item.Item
	err := lib.ReadRequest(r, &item)
	if err != nil {
		lib.WriteResponse(w, err, nil)
		return
	}

	err = h.service.Create(ctx, item)
	if err != nil {
		lib.WriteResponse(w, err, nil)
		return
	}

	lib.WriteResponse(w, nil, nil)
}

func (h *Handler) GetItemHistoryPaginate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request lib.PaginationRequest

	id := r.URL.Query().Get("id")

	request.Page = lib.GetQueryInt(r, "page", 1)
	request.PageSize = lib.GetQueryInt(r, "page_size", 10)

	response, err := h.service.GetItemHistoryPaginate(ctx, id, request)
	if err != nil {
		lib.WriteResponse(w, err, nil)
		return
	}

	lib.WriteResponse(w, nil, response)
}

func (h *Handler) InboundItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type request struct {
		Name string `json:"name"`
		Qty  int    `json:"qty"`
	}

	var req request
	err := lib.ReadRequest(r, &req)
	if err != nil {
		lib.WriteResponse(w, err, nil)
		return
	}

	err = h.service.InboundItem(ctx, req.Name, req.Qty)
	if err != nil {
		lib.WriteResponse(w, err, nil)
		return
	}

	lib.WriteResponse(w, nil, nil)
}

func (h *Handler) OutboundItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type request struct {
		Name string `json:"name"`
		Qty  int    `json:"qty"`
	}

	var req request
	err := lib.ReadRequest(r, &req)
	if err != nil {
		lib.WriteResponse(w, err, nil)
		return
	}

	err = h.service.OutboundItem(ctx, req.Name, req.Qty)
	if err != nil {
		lib.WriteResponse(w, err, nil)
		return
	}

	lib.WriteResponse(w, nil, nil)
}

func (h *Handler) DetectLabels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type request struct {
		ImageBase64 string `json:"image_base64"`
	}

	var req request
	err := lib.ReadRequest(r, &req)
	if err != nil {
		lib.WriteResponse(w, err, nil)
		return
	}

	res, err := h.service.DetectLabels(ctx, req.ImageBase64)
	if err != nil {
		lib.WriteResponse(w, err, nil)
		return
	}

	lib.WriteResponse(w, nil, res)
}
