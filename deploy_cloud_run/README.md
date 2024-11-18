# Project Documentation

## Link google cloud run
- https://temperature-969535875691.us-central1.run.app/78050-328
- https://temperature-969535875691.us-central1.run.app/{CEP}

## Overview

This project is a Go application that provides a web service to get temperature information based on a zip code. It uses the `chi` router for handling HTTP requests and integrates with external APIs to fetch zip code and temperature data.

## Prerequisites

- Go 1.16 or later
- Docker (optional, for containerization)

## Configuration

The application uses environment variables for configuration. Create a `.env` file in the root directory with the following content:

```
PORT=3000
WEATHER_API_KEY=your_weather_api_key
```

Replace `your_weather_api_key` with your actual API key.

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/DiegoOpenheimer/go/deploy_cloud_run.git
    cd deploy_cloud_run
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

## Running the Application

To run the application locally, use the following command:

```sh
go run cmd/main.go
```

The server will start on the port specified in the `.env` file.

## API Endpoints

### Get Temperature

- **URL:** `/ {code}`
- **Method:** `GET`
- **URL Params:**
    - `code` (string): The zip code for which to fetch the temperature.
- **Success Response:**
    - **Code:** 200
    - **Content:** `{"temp_C": "25°C", "temp_F": "25°C", "temp_K": "125°C"}`
- **Error Response:**
    - **Code:** 404
    - **Content:** `{"error": "can not find zipcode"}`
    - **Code:** 422
    - **Content:** `{"error": "invalid zipcode"}`

## Testing

To run the tests, use the following command:

```sh
go test ./...
```

## Docker

To build and run the application using Docker:

1. Build the Docker image:

    ```sh
    docker build -t deploy_cloud_run .
    ```

2. Run the Docker container:

    ```sh
    docker run --env-file .env -p 3000:3000 deploy_cloud_run
    ```