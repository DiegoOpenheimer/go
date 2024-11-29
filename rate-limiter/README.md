# Rate Limiter Project

## Overview

This project implements a rate limiter using Go, Gin, and Redis. The rate limiter restricts the number of requests from a single IP or token within a specified time frame.

## Prerequisites

- Go 1.23.3 or higher
- Redis server
- Git
- Docker and Docker Compose (optional)

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/DiegoOpenheimer/rate-limiter.git
   cd rate-limiter
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

## Configuration

Create a `.env` file in the root directory with the following content:
```
RATE_LIMIT_IP=10
RATE_LIMIT_TOKEN=20
BLOCKED_TIME=60s
REDIS_ADDR=localhost:6379
```

## Running the Application

To start the application, run:
```sh
go run cmd/main.go
```

The server will start on the default port (8080). You can test the rate limiter by sending requests to `http://localhost:8080`.

## Running with Docker Compose

You can also run the application using Docker Compose. This will set up both the application and a Redis server.

1. Build and start the containers:
   ```sh
   docker-compose up --build
   ```

2. The server will be available at `http://localhost:8080`.

## Running Tests

To run the tests, use the following command:
```sh
go test ./...
```

This will execute all unit and integration tests in the project.

## End-to-End Tests

End-to-end tests are located in the `test/e2e` directory. To run these tests, use:
```sh
go test -v ./test/e2e
```

These tests use `miniredis` to simulate a Redis server and validate the rate limiter functionality.

## Project Structure

- `cmd/main.go`: Entry point of the application.
- `config`: Configuration loading.
- `internal/server`: Server and middleware setup.
- `internal/usecases`: Business logic for rate limiting.
- `internal/infra/storage`: Redis storage implementation.
- `test/e2e`: End-to-end tests.