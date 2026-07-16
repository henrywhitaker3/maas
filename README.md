# MaaS (Mutex as a Service)

A small HTTP service that exposes distributed locks backed by Redis. Lock, refresh, and unlock a named subject over a simple JSON API, so multiple processes/services can coordinate without each one embedding its own Redis locking logic.

Locks are implemented with atomic Lua scripts (`SET NX` to acquire, and owner-checked `PEXPIRE`/`DEL` scripts to refresh/release), so refreshing or releasing a lock you don't own is rejected safely even under concurrent access.

## How it works

- **Lock** — `SET key owner NX PX <duration>`. Fails if the key already exists.
- **Refresh** — extends the TTL only if the caller's `owner` matches the current value.
- **Unlock** — deletes the key only if the caller's `owner` matches the current value.

Each lock is identified by a `subject` (the resource being locked) and an `owner` (a UUID identifying who holds it). Only the owner that acquired a lock can refresh or release it.

## API

Base path: `/v1`

| Method | Path        | Description         | Success |
|--------|-------------|----------------------|---------|
| POST   | `/lock`     | Acquire a lock       | `201`   |
| POST   | `/refresh`  | Extend a lock's TTL  | `202`   |
| POST   | `/unlock`   | Release a lock       | `202`   |

### `POST /v1/lock`

```json
{
  "subject": "bongo",
  "owner": "04c43969-3552-4a5b-bd1b-cbcb07aa3d6f",
  "duration": "10s"
}
```

### `POST /v1/refresh`

```json
{
  "subject": "bongo",
  "owner": "04c43969-3552-4a5b-bd1b-cbcb07aa3d6f",
  "duration": "10s"
}
```

### `POST /v1/unlock`

```json
{
  "subject": "bongo",
  "owner": "04c43969-3552-4a5b-bd1b-cbcb07aa3d6f"
}
```

`duration` must be between `1s` and `1h`.

### Error responses

| Status | Meaning                                      |
|--------|-----------------------------------------------|
| `409`  | Lock already exists (on `/lock`), or exists with a different owner (on `/refresh`/`/unlock`) |
| `404`  | Lock not found (on `/refresh`/`/unlock`)       |
| `422`  | Validation error (missing/invalid fields)      |

An OpenAPI spec is generated and served automatically (see `internal/http/http.go`).

## Running locally

Requires [mise](https://mise.jdx.dev/) (see `mise.toml` for pinned Go/Task versions) and Docker.

```bash
task run
```

This brings up a local Redis-compatible store ([Dragonfly](https://github.com/dragonflydb/dragonfly)) via `docker compose` and runs the API with `go run main.go`, listening on port `12345`.

Alternatively, run everything (API + store) in containers:

```bash
docker compose up --build
```

### Configuration

Set via environment variables:

| Variable         | Description              |
|------------------|---------------------------|
| `REDIS_URL`      | Redis/Dragonfly address   |
| `REDIS_PASSWORD` | Redis/Dragonfly password  |

## Go client

A minimal client is available for calling the service from other Go programs:

```go
import "github.com/henrywhitaker3/maas/pkg/client"

c := client.New()

owner := uuid.New()
if err := c.Lock(ctx, "my-subject", owner, 10*time.Second); err != nil {
    // handle: lock already held, etc.
}

// ... do work while holding the lock ...

if err := c.Refresh(ctx, "my-subject", owner, 10*time.Second); err != nil {
    // handle: lock lost/expired, or owned by someone else
}

if err := c.Unlock(ctx, "my-subject", owner); err != nil {
    // handle
}
```

## Project layout

```
main.go                             # entrypoint: wires up Redis client, locker, and HTTP server
internal/locker/                    # Locker interface + Redis implementation (Lua scripts)
internal/http/                      # HTTP server setup, routes, validation, error mapping
internal/http/handlers/locker/      # Lock/Refresh/Unlock HTTP handlers
pkg/locker/                         # Request/response types shared with the client
pkg/client/                         # Go client for calling the API
bruno/MaaS/                         # Bruno API collection for manual testing
```

## Testing requests manually

A [Bruno](https://www.usebruno.com/) collection is included under `bruno/MaaS/` with `Lock` and `Unlock` requests pre-configured against a `{{url}}` environment variable — point it at your local (`http://localhost:12345`) or deployed instance.

## Building

```bash
docker build -t maas .
```

The Dockerfile builds a static Go binary and ships it in a minimal Alpine image, running as `/api`.
