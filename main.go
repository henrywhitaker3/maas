package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/henrywhitaker3/maas/internal/http"
	"github.com/henrywhitaker3/maas/internal/locker"
	"github.com/henrywhitaker3/windowframe/log"
	"github.com/henrywhitaker3/windowframe/redis"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Setup(slog.LevelDebug, os.Stdout)

	client, err := redis.New(ctx, redis.RedisOpts{
		Addr:          os.Getenv("REDIS_URL"),
		Password:      os.Getenv("REDIS_PASSWORD"),
		MaxFlushDelay: time.Microsecond * 100,
	})
	if err != nil {
		panic(err)
	}
	locker := locker.NewRedisLocker(client)

	srv := http.New(locker)
	go func() {
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	_ = srv.Stop(context.Background())
}
