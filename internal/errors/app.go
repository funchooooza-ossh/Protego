package apperrors

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

// ошибки приложения

type AppError error

var (
	//  Клиентские ошибки (4xx)
	ErrNotFound         AppError = errors.New("not found")
	ErrAlreadyExists    AppError = errors.New("already exists")
	ErrInvalidInput     AppError = errors.New("invalid input")
	ErrValidationFailed AppError = errors.New("validation failed")
	ErrUnauthenticated  AppError = errors.New("unauthenticated")
	ErrUnauthorized     AppError = errors.New("unauthorized")
	ErrForbidden        AppError = errors.New("forbidden")
	ErrConflict         AppError = errors.New("conflict")
	ErrRateLimited      AppError = errors.New("rate limited")
	ErrTooManyRequests  AppError = errors.New("too many requests")
	ErrTokenExpired     AppError = errors.New("token expired")
	ErrTokenInvalid     AppError = errors.New("token invalid")

	//  Системные/инфраструктурные (5xx)
	ErrInternal      AppError = errors.New("internal error")
	ErrServiceDown   AppError = errors.New("service unavailable")
	ErrTimeout       AppError = errors.New("operation timeout")
	ErrDependencyErr AppError = errors.New("dependent service failed")
)

var knownErrors = map[error]struct{}{
	ErrNotFound:         {},
	ErrAlreadyExists:    {},
	ErrInvalidInput:     {},
	ErrValidationFailed: {},
	ErrUnauthenticated:  {},
	ErrUnauthorized:     {},
	ErrForbidden:        {},
	ErrConflict:         {},
	ErrRateLimited:      {},
	ErrTooManyRequests:  {},
	ErrTokenExpired:     {},
	ErrTokenInvalid:     {},
	ErrInternal:         {},
	ErrServiceDown:      {},
	ErrTimeout:          {},
	ErrDependencyErr:    {},
}

func IsAppError(err error) bool {
	for known := range knownErrors {
		if errors.Is(err, known) {
			return true
		}
	}
	return false
}

func ToHTTPResponse(err error) (int, string) {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, "resource not found"
	case errors.Is(err, ErrAlreadyExists):
		return http.StatusConflict, "resource already exists"
	case errors.Is(err, ErrInvalidInput),
		errors.Is(err, ErrValidationFailed):
		return http.StatusUnprocessableEntity, "invalid input"
	case errors.Is(err, ErrUnauthenticated):
		return http.StatusUnauthorized, "unauthenticated"
	case errors.Is(err, ErrUnauthorized),
		errors.Is(err, ErrForbidden):
		return http.StatusForbidden, "access denied"
	case errors.Is(err, ErrConflict):
		return http.StatusConflict, "conflict"
	case errors.Is(err, ErrRateLimited),
		errors.Is(err, ErrTooManyRequests):
		return http.StatusTooManyRequests, "too many requests"
	case errors.Is(err, ErrTokenExpired):
		return http.StatusUnauthorized, "token expired"
	case errors.Is(err, ErrTokenInvalid):
		return http.StatusUnauthorized, "token invalid"

	case errors.Is(err, ErrTimeout):
		return http.StatusGatewayTimeout, "request timed out"
	case errors.Is(err, ErrServiceDown),
		errors.Is(err, ErrDependencyErr):
		return http.StatusServiceUnavailable, "service unavailable"

	case errors.Is(err, ErrInternal):
		return http.StatusInternalServerError, "internal error"
	}

	// unknown error
	return http.StatusInternalServerError, "unexpected internal error"
}

type LogLevel string

const (
	Warn  LogLevel = "[WARN]"
	Panic LogLevel = "[PANIC]"
	Fatal LogLevel = "[FATAL]"
	Info  LogLevel = "[INFO]"
)

func LogErr(origin string, err error, level LogLevel) {
	if err == nil {
		return
	}

	log.Printf("%s from: %s | error: %v", level, origin, err)

	switch level {
	case Fatal:
		panic(fmt.Errorf("fatal error from %s: %w", origin, err))
	case Panic:
		panic(err)
	}
}

func ReturnErr(origin string, err error, level LogLevel) error {
	LogErr(origin, err, level)
	return err
}

func BestEffort(origin string, action string, err error) {
	if err != nil {
		LogErr(origin, fmt.Errorf("best-effort %s failed: %w", action, err), Info)
	}
}
