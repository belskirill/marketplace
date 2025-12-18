package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Code string

const (
	CodeInvalidInput       Code = "INVALID_INPUT"       // 400 — неверные данные
	CodeUnauthenticated    Code = "UNAUTHENTICATED"     // 401 — неавторизован
	CodeForbidden          Code = "FORBIDDEN"           // 403 — нет прав
	CodeNotFound           Code = "NOT_FOUND"           // 404 — ресурс не найден
	CodeConflict           Code = "CONFLICT"            // 409 — конфликт данных
	CodeInternal           Code = "INTERNAL"            // 500 — внутренняя ошибка
	CodeServiceUnavailable Code = "SERVICE_UNAVAILABLE" // 503 — временная недоступность
	CodeTimeout            Code = "TIMEOUT"             // 504 — таймаут
)

type HTTPError struct {
	Code    Code              `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
	Err     error             `json:"-"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

func New(code Code, message string, fields map[string]string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
		Fields:  fields,
	}
}

func Wrap(code Code, message string, fields map[string]string, err error) *HTTPError {
	var e *HTTPError
	if errors.As(err, &e) {
		// старый Error сохраняем в Err, новый код и сообщение для клиента
		return &HTTPError{
			Code:    e.Code, // сохраняем код оригинальной ошибки
			Message: fmt.Sprintf("%s: %s", message, e.Message),
			Fields:  fields,
			Err:     e, // внутренняя ошибка
		}
	}

	return &HTTPError{
		Code:    code,
		Message: message,
		Fields:  fields,
		Err:     err,
	}
}

func statusFromCode(code Code) int {
	switch code {
	case CodeInvalidInput:
		return http.StatusBadRequest
	case CodeUnauthenticated:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func Write(w http.ResponseWriter, r *http.Request, op string, logger *zap.Logger, err error) {
	var e *HTTPError
	if !errors.As(err, &e) {
		// если это обычная ошибка, оборачиваем в INTERNAL
		e = &HTTPError{
			Code:    CodeInternal,
			Message: "internal server error",
			Err:     err,
		}
	}

	// Статус для HTTP
	status := statusFromCode(e.Code)

	// Логирование полного traceback, если 500
	if e.Code == CodeInternal {
		logger.Error("request failed",
			zap.String("op", op),
			zap.String("code", string(e.Code)),
			zap.Error(RootCause(err)),
		)

	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"error":      e,
		"request_id": w.Header().Get("X-Request-Id"),
	})
}

func RootCause(err error) error {
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
}

func NewConflict(msg string, fields map[string]string) *HTTPError {
	return New(CodeConflict, msg, fields)
}

func NewInvalidInput(msg string, fields map[string]string) *HTTPError {
	return New(CodeInvalidInput, msg, fields)
}

func NewNotFound(msg string, fields map[string]string) *HTTPError {
	return New(CodeNotFound, msg, fields)
}

func NewUnauthenticated(msg string, fields map[string]string) *HTTPError {
	return New(CodeUnauthenticated, msg, fields)
}
