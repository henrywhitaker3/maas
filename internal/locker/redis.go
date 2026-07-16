package locker

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/henrywhitaker3/windowframe/uuid"
	"github.com/redis/rueidis"
)

var (
	//go:embed refresh.lua
	redisRefresh string
	//go:embed unlock.lua
	redisUnlock string

	refreshScript = rueidis.NewLuaScript(redisRefresh)
	unlockScript  = rueidis.NewLuaScript(redisUnlock)
)

type RedisLocker struct {
	client rueidis.Client
}

func NewRedisLocker(c rueidis.Client) *RedisLocker {
	return &RedisLocker{
		client: c,
	}
}

func (r *RedisLocker) Lock(
	ctx context.Context,
	subject string,
	id uuid.UUID,
	dur time.Duration,
) error {
	cmd := r.client.B().Set().Key(subject).Value(id.String()).Nx().Px(dur).Build()
	if res := r.client.Do(ctx, cmd); res.Error() != nil {
		return fmt.Errorf("%w: %w", ErrLockAlreadyExists, res.Error())
	}
	return nil
}

func (r *RedisLocker) Refresh(
	ctx context.Context,
	subject string,
	owner uuid.UUID,
	duration time.Duration,
) error {
	res := refreshScript.Exec(
		ctx,
		r.client,
		[]string{subject},
		[]string{owner.String(), fmt.Sprint(duration.Milliseconds())},
	)
	if err := res.Error(); err != nil {
		return fmt.Errorf("%w: %w", ErrRefreshError, err)
	}
	n, err := res.ToInt64()
	if err != nil {
		return err
	}
	switch n {
	case 0:
		return ErrRefreshError
	case -1:
		return ErrLockNotFound
	case -2:
		return ErrLockNotOwned
	default:
		return nil
	}
}

func (r *RedisLocker) Unlock(ctx context.Context, subject string, owner uuid.UUID) error {
	resp := unlockScript.Exec(
		ctx,
		r.client,
		[]string{subject},
		[]string{owner.String()},
	)
	if err := resp.Error(); err != nil {
		return fmt.Errorf("%w: %w", ErrUnlock, err)
	}
	n, err := resp.ToInt64()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnlock, err)
	}
	switch n {
	case 0:
		return ErrUnlock
	case -1:
		return ErrLockNotFound
	case -2:
		return ErrLockNotOwned
	default:
		return nil
	}
}

var _ Locker = &RedisLocker{}
