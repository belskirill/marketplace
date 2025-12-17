package dto

import (
	httperr "bots/pkg/errors"
)

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"omitempty,min=3,max=20"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (r *LoginUserRequest) Validate() error {
	if r.Username == "" && r.Email == "" {
		return httperr.NewInvalidInput("either username or email must be provided", map[string]string{
			"username": "required if email is empty",
			"email":    "required if username is empty",
		})
	}
	return nil
}
