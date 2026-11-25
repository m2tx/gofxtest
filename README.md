# GoFxTest

GoFxTest is a starter project built with [Go](https://golang.org/) and [Uber FX](https://github.com/uber-go/fx), designed to help you quickly bootstrap robust, modular Go applications with dependency injection.

## Features

- Modular architecture using Uber FX
- REST API with [Gin](https://github.com/gin-gonic/gin)
- MongoDB integration
- RabbitMQ pub/sub
- OpenTelemetry metrics exported in Prometheus format on `/metrics`
- Example `events` CRUD resource showing how the layers fit together
- Swagger API documentation
- Makefile for common development tasks
- Linting, auditing, and test coverage support

## Getting Started

### Prerequisites

- Go 1.27+
- Docker with Docker Compose, which the Makefile uses to run MongoDB, RabbitMQ, the tests, the linter and the auditors

### Installation

Install project dependencies:

```shell
make install
```

### Configuration

Configuration is read from environment variables. `make run` and the VS Code launch configuration load them from a `.env` file in the project root, which is gitignored.

| Variable                | Default                              | Description                           |
|-------------------------|--------------------------------------|---------------------------------------|
| `NAME`                  | `gofxtest`                           | Service name used by the logger       |
| `HTTP_PORT`             | `8080`                               | HTTP server port                      |
| `HTTP_READ_TIMEOUT`     | `10s`                                | HTTP read timeout                     |
| `HTTP_WRITE_TIMEOUT`    | `10s`                                | HTTP write timeout                    |
| `HTTP_SHUTDOWN_TIMEOUT` | `5s`                                 | Graceful shutdown timeout             |
| `MONGO_URL`             | required                             | MongoDB connection URL                |
| `MONGO_DATABASE`        | required                             | MongoDB database name                 |
| `RABBITMQ_URL`          | `amqp://guest:guest@localhost:5672/` | RabbitMQ connection URL               |
| `RABBITMQ_QUEUE_PREFIX` | `gofxtest`                           | Prefix for the subscriber queue names |

Example `.env` for the local Docker environment:

```shell
MONGO_URL=mongodb://localhost:27017
MONGO_DATABASE=gofxtest
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```

### Running the Application

Start MongoDB and RabbitMQ, then run the server with the variables from `.env`:

```shell
make run
```

Or, with the environment already up and the variables exported:

```shell
go run ./cmd/server/main.go
```

## Example Implementation: Events

The `events` resource is a sample feature, not part of the template's core. It shows the pattern to copy for your own resources:

| Layer                     | Location                         | Role                                              |
|---------------------------|----------------------------------|---------------------------------------------------|
| Domain model + interfaces | `domain/event.go`                | `Event` entity, `EventService`, `EventRepository` |
| Service                   | `domain/event/`                  | Business logic; publishes to RabbitMQ on changes  |
| Repository                | `internal/repository/mongo/`     | MongoDB persistence                               |
| HTTP routes               | `internal/http/event.go`         | Gin handlers with Swagger annotations             |
| Wiring                    | `cmd/server/main.go`             | `fx.Provide` registrations and queue subscribers  |

It exposes these routes:

| Method   | Path           | Description        |
|----------|----------------|--------------------|
| `GET`    | `/events`      | List all events    |
| `GET`    | `/events/{ID}` | Get an event by ID |
| `POST`   | `/events`      | Create an event    |
| `PUT`    | `/events/{ID}` | Update an event    |
| `DELETE` | `/events/{ID}` | Delete an event    |

Each create, update and delete publishes the event to the RabbitMQ topic `event-created`, `event-updated` or `event-deleted`. The subscribers in `cmd/server/main.go` only log the messages.

To add your own resource, follow the same layers and register its constructors in `cmd/server/main.go`. Register HTTP routes with `asRouteHandler(...)`. To start a project without the example, remove those files and their registrations.

The template's own endpoints are `GET /health` (health check), `GET /metrics` (Prometheus metrics) and `/swagger/*` (API docs).

## Metrics

Metrics are collected with the [OpenTelemetry](https://opentelemetry.io/) SDK and exposed on `GET /metrics` in Prometheus format:

| Metric                                  | Source                                                            |
|-----------------------------------------|-------------------------------------------------------------------|
| `http_server_request_duration_seconds`  | Every HTTP request, labelled by route, method and status (`otelgin`) |
| `db_client_operation_duration_seconds`  | Every MongoDB command, labelled by operation and status           |
| `messaging_client_published_messages_total` | Messages published to RabbitMQ, labelled by topic and status  |
| `messaging_process_duration_seconds`    | Duration of the subscriber handlers, labelled by topic            |
| `go_*`, `process_*`                     | Go runtime and process metrics                                    |

`target_info` carries the service name from `NAME`. To add your own metrics, inject `metric.MeterProvider` into a constructor and create instruments from `meterProvider.Meter("<package path>")`.

## Makefile Commands

| Command               | Description                                              |
|-----------------------|----------------------------------------------------------|
| `make install`        | Install Go dependencies                                  |
| `make install-tools`  | Install the swag CLI                                     |
| `make env-up`         | Start MongoDB and RabbitMQ                               |
| `make env-down`       | Stop and remove MongoDB and RabbitMQ                     |
| `make run`            | Start the environment and run the server                 |
| `make build`          | Build the binary in Docker                               |
| `make image`          | Build the runtime Docker image                           |
| `make lint`           | Run golangci-lint                                        |
| `make audit`          | Audit dependencies for vulnerabilities with OSV-Scanner  |
| `make audit-gosec`    | Run the gosec security scanner                           |
| `make test`           | Run tests against MongoDB and write reports to `coverage/` |
| `make coverage`       | Run tests and open the HTML coverage report              |
| `make swag`           | Generate Swagger API docs                                |
| `make destroy`        | Remove the test container and the environment            |

## API Documentation

Generate Swagger documentation after changing routes or their annotations:

```shell
make swag
```

With the server running, the Swagger UI is at `http://localhost:8080/swagger/index.html`.
