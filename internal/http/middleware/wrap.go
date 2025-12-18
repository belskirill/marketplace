package middleware

import (
	httperr "bots/pkg/errors"
	"net/http"

	"go.uber.org/zap"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func Wrap(h AppHandler, logger *zap.Logger, op string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			httperr.Write(w, r, op, logger, err)
		}
	}
}
