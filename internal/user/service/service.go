package service

import (
	"bots/internal/security"
	"bots/internal/user/http/dto"
	"bots/internal/user/repositories/postgres"
	httperr "bots/pkg/errors"
	"context"

	"go.uber.org/zap"
)

type UserService interface {
	RegisterUser(ctx context.Context, req dto.RegisterUserRequest) (dto.UserResponse, error)
}

func NewService(db postgres.DatabaseRepository, logger *zap.Logger) UserService {
	return &Service{
		db:     db,
		logger: logger,
	}
}

type Service struct {
	db     postgres.DatabaseRepository
	logger *zap.Logger
}

func (s *Service) RegisterUser(ctx context.Context, req dto.RegisterUserRequest) (dto.UserResponse, error) {
	hashPass, err := security.HashPassword(req.Password)
	if err != nil {
		return dto.UserResponse{}, httperr.New(httperr.CodeInternal,
			"failed to hash password",
			nil,
		)
	}
	req.Password = hashPass
	resp, err := s.db.Create(ctx, &req)
	if err != nil {
		switch err {
		case postgres.ErrEmailTaken:
			return dto.UserResponse{}, httperr.NewConflict("Пользователь с таким email уже существует",
				map[string]string{"email": req.Email})
		case postgres.ErrUsernameTaken:
			return dto.UserResponse{}, httperr.NewConflict("Пользователь с таким username уже существует",
				map[string]string{"username": req.Username})
		default:
			return dto.UserResponse{}, httperr.Wrap(httperr.CodeInternal, "Не удалось создать пользователя", nil, err)
		}
	}

	return resp, nil
}
