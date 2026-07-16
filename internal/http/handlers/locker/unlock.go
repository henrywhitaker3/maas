// Package locker
package locker

import (
	"net/http"

	intLocker "github.com/henrywhitaker3/maas/internal/locker"
	"github.com/henrywhitaker3/maas/pkg/locker"
	"github.com/henrywhitaker3/windowframe/http/common"
	"github.com/labstack/echo/v4"
)

type UnlockHandler struct {
	locker intLocker.Locker
}

func NewUnlockHandler(l intLocker.Locker) *UnlockHandler {
	return &UnlockHandler{locker: l}
}

func (l *UnlockHandler) Handler() common.Handler[locker.UnlockRequest, any] {
	return func(c echo.Context, req locker.UnlockRequest) (*any, error) {
		if err := l.locker.Unlock(
			c.Request().Context(),
			req.Subject,
			req.Owner,
		); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func (l *UnlockHandler) Metadata() common.Metadata {
	return common.Metadata{
		Name:         "Unlock",
		Description:  "Unlock a key",
		Tag:          "Locks",
		Code:         http.StatusAccepted,
		Method:       http.MethodPost,
		Path:         "/v1/unlock",
		GenerateSpec: true,
	}
}

func (l *UnlockHandler) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
