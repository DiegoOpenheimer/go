# OpenTelemetry Microservices Project

This project is a study on implementing distributed tracing using OpenTelemetry in a microservices architecture. The project consists of two services (`service_a` and `service_b`) and uses Zipkin for tracing and the OpenTelemetry Collector for exporting trace data.

## Project Structure

- `service_a`: The first microservice that handles incoming requests and forwards them to `service_b`.
- `service_b`: The second microservice that processes requests from `service_a`.
- `otel-collector`: The OpenTelemetry Collector that collects and exports trace data.
- `zipkin`: The Zipkin server for visualizing trace data.

## Prerequisites

- Docker
- Docker Compose
- Go 1.18+

## Configuration

### Environment Variables

The project uses environment variables to configure the services. You can set these variables in a `.env` file or directly in the `docker-compose.yaml` file.

- `SERVICE_A_PORT`: The port on which `service_a` runs (default: 3000).
- `SERVICE_B_PORT`: The port on which `service_b` runs (default: 3001).
- `SERVICE_B_URL`: The URL of `service_b` (default: `http://service_b:3001`).
- `WEATHER_API_KEY`: The API key for the weather service used by `service_b`.
- `SERVICE_NAME`: The name of the service for OpenTelemetry.
- `OTEL_EXPORTER_OTLP_ENDPOINT`: The endpoint for the OpenTelemetry Collector (default: `otel-collector:4317`).

### Docker Compose

The `docker-compose.yaml` file defines the services and their dependencies. To start the services, run:

```sh
docker-compose up --build