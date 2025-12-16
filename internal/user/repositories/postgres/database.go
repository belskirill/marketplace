package postgres

import (
	"bots/internal/user/http/dto"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	ErrEmailTaken    = errors.New("email already exists")
	ErrUsernameTaken = errors.New("username already exists")
)

type DatabaseRepository interface {
	GetUsers(ctx context.Context)
	Create(ctx context.Context, req *dto.RegisterUserRequest) (dto.UserResponse, error)
}

func NewDatabaseRepository(db *sql.DB, logger *zap.Logger) DatabaseRepository {
	return &Database{
		db:     db,
		logger: logger,
	}
}

type Database struct {
	db     *sql.DB
	logger *zap.Logger
}

func (repo *Database) GetUsers(ctx context.Context) {
	//dsdsd
}

func (repo *Database) Create(ctx context.Context, req *dto.RegisterUserRequest) (dto.UserResponse, error) {
	var resp dto.UserResponse
	const query = `
        INSERT INTO users (username, email, password_hash)
        VALUES ($1, $2, $3)
        RETURNING id, username, email
    `
	err := repo.db.QueryRowContext(ctx, query,
		req.Username, req.Email, req.Password,
	).Scan(&resp.ID, &resp.Username, &resp.Email)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && string(pqErr.Code) == "23505" {
			switch pqErr.Constraint {
			case "users_email_key":
				return dto.UserResponse{}, ErrEmailTaken
			case "users_username_key":
				return dto.UserResponse{}, ErrUsernameTaken
			}
		}
		return dto.UserResponse{}, fmt.Errorf("postgres insert users: %w", err)
	}

	return resp, nil
}
