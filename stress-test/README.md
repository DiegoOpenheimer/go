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

Clone the repository and build the application

## Usage

Run the application with the following arguments:

```bash
./stress-test --url=<API_URL> --requests=<NUMBER_OF_REQUESTS> --concurrency=<CONCURRENCY_LEVEL>
```

Example:

```bash
./stress-test --url=<API_URL> --requests=100 --concurrency=10
```

## Configuration

The application uses the following environment variables for configuration:

- `URL`: The API endpoint to be tested.
- `REQUESTS`: Total number of requests to perform.
- `CONCURRENCY`: Number of concurrent workers.

Set the environment variables:

```bash
export URL=<API_URL>
export REQUESTS=100
export CONCURRENCY=10
```

Alternatively, you can pass these values directly as command-line arguments.

## Running with Docker

You can also run the application using Docker. The project is available as a Docker image hosted on Docker Hub.

### Pull the image

You can build the dockerfile or Pull the Docker image from Docker Hub:

```bash
docker pull diegoopenheimer/stress-test
```

### Run the application

Run the application using the Docker image:

```bash
docker run --rm --name stress-test diegoopenheimer/stress-test --url=https://www.google.com --requests=100 --concurrency=10
```

The `--rm` flag ensures the container is removed after execution. Replace the `--url`, `--requests`, and `--concurrency` values as needed.

## Validation Errors

The application validates the input parameters and provides meaningful error messages. Examples:

- **Missing URL**: `Error: URL is required.`
- **Invalid URL**: `Error: URL must be a valid URL.`

