# GoFxTest

GoFxTest is a starter project built with [Go](https://golang.org/) and [Uber FX](https://github.com/uber-go/fx), designed to help you quickly bootstrap robust, modular Go applications with dependency injection.

## Features

- Modular architecture using Uber FX
- MongoDB integration
- Swagger API documentation
- Makefile for common development tasks
- Linting, auditing, and test coverage support

## Getting Started

### Prerequisites

- Go 1.20+ installed
- MongoDB instance running

### Installation

Install project dependencies:

```shell
make install
```

### Running the Application

```shell
go run ./cmd/server/main.go
```

## Makefile Commands

| Command      | Description                  |
|--------------|-----------------------------|
| `make install` | Install Go dependencies      |
| `make lint`    | Run code linter              |
| `make audit`   | Audit dependencies for vulnerabilities |
| `make test`    | Run tests and show coverage  |
| `make swag`    | Generate Swagger API docs    |

## API Documentation

Generate Swagger documentation:

```shell
make swag
```