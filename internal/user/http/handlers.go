package http

import (
	"bots/internal/user/http/dto"
	"bots/internal/user/service"
	httperr "bots/pkg/errors"
	"encoding/json"
	"fmt"
	"net/http"

	mdlwr "bots/internal/http"
	"github.com/go-playground/validator/v10"
)

type Handler interface {
	CreateUser(w http.ResponseWriter, r *http.Request) error
	Login(w http.ResponseWriter, r *http.Request) error
}

type UserHandler struct {
	service  service.UserService
	validate *validator.Validate
}

func NewHandler(srv service.UserService) Handler {
	return &UserHandler{
		service:  srv,
		validate: validator.New(),
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	req := dto.RegisterUserRequest{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := dec.Decode(&req); err != nil {
		return httperr.NewInvalidInput("invalid JSON", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		fields := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			fields[e.Field()] = fmt.Sprintf("failed on %s", e.Tag())
		}
		return httperr.NewInvalidInput("validation failed", fields)
	}
	resp, err := h.service.RegisterUser(r.Context(), req)
	if err != nil {
		return err
	}
	return mdlwr.RespondJSON(w, http.StatusCreated, resp)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) error {
	req := dto.LoginUserRequest{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		return err
	}

	if err := h.validate.Struct(req); err != nil {
		fields := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Println(e.Field())
			fields[e.Field()] = fmt.Sprintf("failed on %s", e.Tag())
		}
		return httperr.NewInvalidInput("validation failed", fields)
	}

	if err := req.Validate(); err != nil {
		return err
	}
	res, err := h.service.LoginUser(r.Context(), req)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    res,
		Path:     "/",
		HttpOnly: true,                 // JS не видит токен
		SameSite: http.SameSiteLaxMode, // или Strict
		MaxAge:   15 * 60,              // 15 минут
	})
	return mdlwr.RespondJSON(w, http.StatusNoContent, nil)
}
