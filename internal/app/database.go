package app

import (
	"context"
	"database/sql"
	"time"
)

func Connect(ctx context.Context, dsn string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := conn.PingContext(pingCtx); err != nil {
		return nil, err
	}

	return conn, nil
}
