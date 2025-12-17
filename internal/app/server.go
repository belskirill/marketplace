package app

import (
	hndlr "bots/internal/user/http"
	"bots/internal/user/repositories/postgres"
	"bots/internal/user/service"
	"context"
	"database/sql"
	"log"
	"net/http"

	"go.uber.org/zap"
)

func NewServer(addr string) (*http.Server, *sql.DB) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	mux := http.NewServeMux()

	ctx := context.Background()

	cfg, err := Load()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := Connect(ctx, cfg.DB.DSN())
	if err != nil {
		log.Fatal(err)
	}
	dbRepo := postgres.NewDatabaseRepository(conn, logger)
	userSrv := service.NewService(dbRepo, logger)
	userHandler := hndlr.NewHandler(userSrv, logger)

	mux.HandleFunc("/register", userHandler.CreateUser)
	mux.HandleFunc("/login", userHandler.Login)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return srv, conn
}
