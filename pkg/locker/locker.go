// Package locker
package locker

import (
	"github.com/henrywhitaker3/windowframe/duration"
	"github.com/henrywhitaker3/windowframe/uuid"
)

type LockRequest struct {
	Subject  string                  `json:"subject"  validate:"required"               required:"true"`
	Owner    uuid.UUID               `json:"owner"    validate:"required"               required:"true"`
	Duration duration.StringDuration `json:"duration" validate:"required,lock_duration" required:"true"`
}

type RefreshRequest struct {
	Subject  string                  `json:"subject"  validate:"required"               required:"true"`
	Owner    uuid.UUID               `json:"owner"    validate:"required"               required:"true"`
	Duration duration.StringDuration `json:"duration" validate:"required,lock_duration" required:"true"`
}

type UnlockRequest struct {
	Subject string    `json:"subject" validate:"required" required:"true"`
	Owner   uuid.UUID `json:"owner"   validate:"required" required:"true"`
}
