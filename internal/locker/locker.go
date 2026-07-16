// Package locker
package locker

import (
	"context"
	"errors"
	"time"

	"github.com/henrywhitaker3/windowframe/uuid"
)

type Locker interface {
	// Lock a key for a certian duration
	Lock(ctx context.Context, subject string, owner uuid.UUID, duration time.Duration) error
	// Refresh a lock for a duration
	Refresh(ctx context.Context, subject string, owner uuid.UUID, duration time.Duration) error
	// Release a lock
	Unlock(ctx context.Context, subject string, owner uuid.UUID) error
}

var (
	ErrLockAlreadyExists = errors.New("lock already exists")
	ErrRefreshError      = errors.New("refresh error")
	ErrUnlock            = errors.New("unlock")
	ErrLockNotFound      = errors.New("lock not found")
	ErrLockNotOwned      = errors.New("lock has different owner")
)
