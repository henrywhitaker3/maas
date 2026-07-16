// Package http
package http

import (
	"errors"
	"log/slog"
	"time"

	stdhttp "net/http"

	"github.com/go-playground/validator/v10"
	"github.com/henrywhitaker3/maas/internal/http/handlers/locker"
	intLocker "github.com/henrywhitaker3/maas/internal/locker"
	"github.com/henrywhitaker3/windowframe/duration"
	"github.com/henrywhitaker3/windowframe/http"
)

func New(l intLocker.Locker) *http.HTTP {
	h := http.New(http.HTTPOpts{
		Port: 12345,
		Openapi: http.OpenapiOpts{
			Enabled:        true,
			ServiceName:    "MaaS",
			ServiceVersion: "1.0.0",
			PublicURL:      "https://mass.henrywhitaker.com",
		},
		Logger: slog.With("component", "http"),
	})
	registerRoutes(h, l)
	registerValidations(h)
	handleErrors(h)
	return h
}

func handleErrors(h *http.HTTP) {
	h.HandleErrors(
		func(err error) (int, any, bool) {
			if errors.Is(err, intLocker.ErrLockAlreadyExists) {
				return stdhttp.StatusConflict, http.NewError("lock exists"), true
			}
			return 0, nil, false
		},
		func(err error) (int, any, bool) {
			if errors.Is(err, intLocker.ErrLockNotFound) {
				return stdhttp.StatusNotFound, http.NewError("lock not found"), true
			}
			return 0, nil, false
		},
		func(err error) (int, any, bool) {
			if errors.Is(err, intLocker.ErrLockNotOwned) {
				return stdhttp.StatusConflict, http.NewError("lock has different owner"), true
			}
			return 0, nil, false
		},
	)
}

func registerRoutes(h *http.HTTP, l intLocker.Locker) {
	http.Register(h, locker.NewLockHandler(l))
	http.Register(h, locker.NewUnlockHandler(l))
	http.Register(h, locker.NewRefreshHandler(l))
}

func registerValidations(h *http.HTTP) {
	h.Validator.RegisterValidation("lock_duration", func(fl validator.FieldLevel) bool {
		strDur, ok := fl.Field().Interface().(duration.StringDuration)
		if !ok {
			return false
		}
		dur := strDur.Duration()

		if dur >= time.Second && dur <= time.Hour {
			return true
		}

		return false
	})
}
