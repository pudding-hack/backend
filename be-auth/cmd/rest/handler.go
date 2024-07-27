package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pudding-hack/backend/be-auth/internal/use_case"
	"github.com/pudding-hack/backend/lib"
)

type service interface {
	Login(ctx context.Context, username, password string) (res use_case.LoginResponse, err error)
	GetCurrentUser(ctx context.Context, accessToken string) (res use_case.User, err error)
	RefreshToken(ctx context.Context, refreshToken string) (res use_case.LoginResponse, err error)
	Logout(ctx context.Context, accessToken string) (err error)
}

type Handler struct {
	service service
}

func NewHandler(service service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.WriteResponse(w, lib.NewErrBadRequest(err.Error()), nil)
		return
	}

	res, err := h.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		lib.WriteResponse(w, err, nil)
		return
	}

	lib.WriteResponse(w, nil, res)
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")
	splitAccessToken := strings.Split(accessToken, "Bearer ")
	if len(splitAccessToken) != 2 {
		lib.WriteResponse(w, lib.NewErrBadRequest("invalid access token"), nil)
		return
	}

	res, err := h.service.GetCurrentUser(r.Context(), splitAccessToken[1])
	if err != nil {
		lib.WriteResponse(w, err, nil)
		return
	}

	lib.WriteResponse(w, nil, res)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("Authorization")
	splitRefreshToken := strings.Split(refreshToken, "Bearer ")

	if len(splitRefreshToken) != 2 {
		lib.WriteResponse(w, lib.NewErrBadRequest("invalid refresh token"), nil)
		return
	}

	res, err := h.service.RefreshToken(r.Context(), splitRefreshToken[1])
	if err != nil {
		lib.WriteResponse(w, err, nil)
		return
	}

	lib.WriteResponse(w, nil, res)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")
	splitAccessToken := strings.Split(accessToken, "Bearer ")
	if len(splitAccessToken) != 2 {
		lib.WriteResponse(w, lib.NewErrBadRequest("invalid access token"), nil)
		return
	}

	err := h.service.Logout(r.Context(), splitAccessToken[1])
	if err != nil {
		lib.WriteResponse(w, err, nil)
		return
	}

	lib.WriteResponse(w, nil, nil)
}
