package service

import (
	httperr "bots/internal/errors"
	"bots/internal/user/http/dto"
	"bots/internal/user/repositories/postgres"
	"context"

	"go.uber.org/zap"
)

type UserService interface {
	GetUsers(ctx context.Context, req dto.RegisterUserRequest) (dto.UserResponse, error)
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

func (s *Service) GetUsers(ctx context.Context, req dto.RegisterUserRequest) (dto.UserResponse, error) {
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
