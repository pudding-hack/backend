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
	GetByID(ctx context.Context, id string) (use_case.Item, error)
	Create(ctx context.Context, item item.Item) error
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
	id := r.URL.Query().Get("id")

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
