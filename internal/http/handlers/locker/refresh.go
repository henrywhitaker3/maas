// Package locker
package locker

import (
	"net/http"

	intLocker "github.com/henrywhitaker3/maas/internal/locker"
	"github.com/henrywhitaker3/maas/pkg/locker"
	"github.com/henrywhitaker3/windowframe/http/common"
	"github.com/labstack/echo/v4"
)

type RefreshHandler struct {
	locker intLocker.Locker
}

func NewRefreshHandler(l intLocker.Locker) *RefreshHandler {
	return &RefreshHandler{locker: l}
}

func (l *RefreshHandler) Handler() common.Handler[locker.RefreshRequest, any] {
	return func(c echo.Context, req locker.RefreshRequest) (*any, error) {
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

func (l *RefreshHandler) Metadata() common.Metadata {
	return common.Metadata{
		Name:         "Refresh",
		Description:  "Refresh a lock",
		Tag:          "Locks",
		Code:         http.StatusAccepted,
		Method:       http.MethodPost,
		Path:         "/v1/refresh",
		GenerateSpec: true,
	}
}

func (l *RefreshHandler) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
