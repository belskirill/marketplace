package app

import (
	wrap "bots/internal/http/middleware"
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
	userHandler := hndlr.NewHandler(userSrv)

	mux.HandleFunc("/register", wrap.Wrap(userHandler.CreateUser, logger, "user.register"))
	mux.HandleFunc("/login", wrap.Wrap(userHandler.Login, logger, "user.login"))

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return srv, conn
}
