# Stress Test CLI

A command-line tool written in Go for stress testing APIs. This project allows you to simulate multiple concurrent requests to test the performance and reliability of a given API.

## Features

- **Configurable requests and concurrency**: Set the number of requests and concurrency level via command-line arguments.
- **Input validation**: Ensures that input parameters such as the URL are valid.
- **Built-in error handling**: Provides meaningful error messages for invalid configurations.

## Requirements

- Go 1.23.3 or higher
- Environment variables configured for URL, concurrency, and requests.

## Installation

Clone the repository and build the application:

## Usage

Run the application with the following arguments:

```bash
./stress-test --url=<API_URL> --requests=<NUMBER_OF_REQUESTS> --concurrency=<CONCURRENCY_LEVEL>
```

Example:

```bash
./stress-test --url=https://viacep.com.br/ws/66813880/json --requests=100 --concurrency=10
```

## Configuration

The application uses the following environment variables for configuration:

- `URL`: The API endpoint to be tested.
- `REQUESTS`: Total number of requests to perform.
- `CONCURRENCY`: Number of concurrent workers.

Set the environment variables:

```bash
export URL=https://viacep.com.br/ws/66813880/json
export REQUESTS=100
export CONCURRENCY=10
```

Alternatively, you can pass these values directly as command-line arguments.

## Validation Errors

The application validates the input parameters and provides meaningful error messages. Examples:

- **Missing URL**: `Error: URL is required.`
- **Invalid URL**: `Error: URL must be a valid URL.`

## Development

### Dependencies

Install dependencies with:

```bash
go mod tidy
```
