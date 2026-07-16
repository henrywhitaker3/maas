// Package locker
package locker

import (
	"net/http"

	intLocker "github.com/henrywhitaker3/maas/internal/locker"
	"github.com/henrywhitaker3/maas/pkg/locker"
	"github.com/henrywhitaker3/windowframe/http/common"
	"github.com/labstack/echo/v4"
)

type LockHandler struct {
	locker intLocker.Locker
}

func NewLockHandler(l intLocker.Locker) *LockHandler {
	return &LockHandler{locker: l}
}

func (l *LockHandler) Handler() common.Handler[locker.LockRequest, any] {
	return func(c echo.Context, req locker.LockRequest) (*any, error) {
		if err := l.locker.Lock(
			c.Request().Context(),
			req.Subject,
			req.Owner,
			req.Duration.Duration(),
		); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func (l *LockHandler) Metadata() common.Metadata {
	return common.Metadata{
		Name:         "Lock",
		Description:  "Lock a key",
		Tag:          "Locks",
		Code:         http.StatusCreated,
		Method:       http.MethodPost,
		Path:         "/v1/lock",
		GenerateSpec: true,
	}
}

func (l *LockHandler) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
