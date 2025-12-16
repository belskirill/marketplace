package http

import (
	httperr "bots/internal/errors"
	"bots/internal/user/http/dto"
	"bots/internal/user/service"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Handler interface {
	GetUsers(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	service  service.UserService
	validate *validator.Validate
	logger   *zap.Logger
}

func NewHandler(srv service.UserService, logger *zap.Logger) Handler {
	return &UserHandler{
		service:  srv,
		validate: validator.New(),
		logger:   logger,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	req := dto.RegisterUserRequest{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		httperr.Write(w, r, "users.register", h.logger,
			httperr.NewInvalidInput("invalid JSON", nil))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httperr.Write(w, r, "users.register", h.logger,
			httperr.NewInvalidInput("invalid email", nil))
		return
	}
	resp, err := h.service.GetUsers(r.Context(), req)
	if err != nil {
		httperr.Write(w, r, "users.register", h.logger, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// Если вдруг ошибка при кодировании
		httperr.Write(w, r, "users.register", h.logger, httperr.Wrap(httperr.CodeInternal, "failed to encode response", nil, err))
	}
}

type S struct {
}
