# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Go 1.25 starter service built on Uber FX (dependency injection) with a Gin HTTP API, MongoDB persistence, and RabbitMQ pub/sub. Module path: `github.com/m2tx/gofxtest`.

## Commands

Most Makefile targets run inside Docker (see `build/`), so Docker is required.

- `make env-up` / `make env-down`: start/stop MongoDB (27017) and RabbitMQ (5672, management UI on 15672) via `build/docker-compose.yml`. `make env-up-mongo` starts only Mongo.
- `make run`: brings up the environment, loads `.env` (gitignored; keys: `NAME`, `HTTP_PORT`, `MONGO_URL`, `MONGO_DATABASE`, `RABBITMQ_URL`) and runs `go run ./cmd/server/main.go`.
- `make build`: builds the binary in the Dockerfile `build` stage. `make image` builds the runtime image.
- `make test`: starts Mongo, runs `go test ./...` inside the Dockerfile `test` stage on the compose network, and copies coverage into `coverage/` (mock files are excluded from coverage). `make coverage` also opens the HTML report.
- `make lint`: golangci-lint v2 in Docker (`.golangci.yml`, `standard` linters, test files not linted).
- `make audit` (osv-scanner), `make audit-gosec`.
- `make swag`: regenerates `docs/` from the swag annotations on handlers in `internal/http`. Run this after changing routes or annotations. Install the CLI with `make install-tools`.

Running tests locally without Docker:

```shell
make env-up-mongo
MONGO_URL=mongodb://localhost:27017 MONGO_DATABASE=test go test ./...
go test ./domain/event -run TestEventService_Get -v       # single test (no infra needed)
```

`internal/repository/mongo` tests are integration tests. Their `TestMain` connects to a real MongoDB using `MONGO_URL`/`MONGO_DATABASE` and panics if those are unset. The `domain/event` and `internal/env` tests are pure unit tests.

CI (`.github/workflows/pipeline.yml`) runs Build → Lint → Audit (non-blocking) → Test on PRs, then posts a coverage badge comment. A coverage figure of 80% or less is shown in red.

## Architecture

The `events` resource (domain, service, Mongo repository, HTTP route and queue subscribers) is an example implementation that shows the pattern. It is not core template functionality. New resources should follow the same layering.

**Wiring lives entirely in `cmd/server/main.go`.** Every component is an `fx.Provide` constructor, and lifecycle concerns (starting the HTTP server, Mongo disconnect, queue subscriptions) are `fx.Invoke` functions that register `fx.Lifecycle` hooks. To add a component, write a constructor that takes its dependencies as parameters and register it there.

- **Config**: each component declares its own config struct with `env:"VAR"`, `default:"..."`, `required:"true"`, and `min`/`max` tags. It is provided via `fx.Provide(env.New[ConfigType])`. `internal/env` is a custom reflection-based loader that supports string, int, float, bool, and `time.Duration`. Other field types return an error.
- **Layers**:
  - `domain/` holds entities and the `EventService` / `EventRepository` interfaces.
  - `domain/event/` implements the service. Each mutation persists through the repository, then publishes to a topic (`event-created` / `event-updated` / `event-deleted`).
  - `internal/repository/mongo/` implements the repositories. IDs are stored as hex-string `_id`s, not `ObjectID`s.
  - `internal/http/` holds the transport.
  - `internal/queue/` holds the `Publisher`/`Subscriber` interfaces and the RabbitMQ implementation.
- **HTTP routes use an FX value group.** Each route type implements `http.RouteHandler` (`Register(*gin.Engine)`) and is provided through the `asRouteHandler(...)` helper, which tags it into `group:"routeHandlers"`. `http.NewHandler` receives the whole group and registers each route. A new endpoint therefore needs its own `RouteHandler` constructor plus one `fx.Provide(asRouteHandler(...))` line. Do not edit the handler to add it.
- **Queue**: the single RabbitMQ instance is exposed as both `queue.Publisher` and `queue.Subscriber` via `asQueue(...)`. Topics map to durable fanout exchanges, and each subscriber gets an exclusive queue named `<RABBITMQ_QUEUE_PREFIX>-<topic>`. Exchanges are declared only in `Subscribe`, so publishing to a topic nobody has subscribed to fails. Non-`queue.Message` values passed to `Publish` are JSON-encoded.

## Testing conventions

- Tests use `testify/assert` and the external `_test` package (e.g. `package event_test`).
- Mocks are hand-written next to their implementations (`*_mock.go`, e.g. `mongo.EventRepositoryMock`, `queue.QueueMock`). They are structs with a `XxxFn` func field and an `XxxCount` call counter per method. Set only the funcs a test needs, because calling an unset one panics. Assert on the counters.
- Several `internal/http/*_test.go` files are empty placeholders (package declaration only).
