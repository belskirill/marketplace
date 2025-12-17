package http

import (
	"bots/internal/user/http/dto"
	"bots/internal/user/service"
	httperr "bots/pkg/errors"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Handler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
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

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	req := dto.RegisterUserRequest{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := dec.Decode(&req); err != nil {
		httperr.Write(w, r, "users.register", h.logger,
			httperr.NewInvalidInput("invalid JSON", nil))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		fields := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			fields[e.Field()] = fmt.Sprintf("failed on %s", e.Tag())
		}
		httperr.Write(w, r, "users.register", h.logger,
			httperr.NewInvalidInput("validation failed", fields))
		return
	}
	resp, err := h.service.RegisterUser(r.Context(), req)
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
	h.logger.Info("user created",
		zap.String("username", req.Username),
		zap.String("email", req.Email),
	)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	req := dto.LoginUserRequest{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		httperr.Write(w, r, "user.login", h.logger,
			httperr.NewInvalidInput("invalid JSON", nil))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		fields := make(map[string]string)
		fmt.Println(err.(validator.ValidationErrors))
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Println(e.Field())
			fields[e.Field()] = fmt.Sprintf("failed on %s", e.Tag())
		}
		httperr.Write(w, r, "users.login", h.logger,
			httperr.NewInvalidInput("validation failed", fields))
		return
	}

	if err := req.Validate(); err != nil {
		httperr.Write(w, r, "users.login", h.logger, err)
		return
	}
}
