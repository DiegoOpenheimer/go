# Auction Project

## Description
This study project is an auction system developed in Go. It allows creating auctions, placing bids, and automatically closing auctions after a time interval, this implementation is in create_auction_usecase.

## How to Run

### Prerequisites
- Go 1.23 or higher
- MongoDB

### Steps to Run

1. Install dependencies:
   ```sh
   go mod tidy
   ```

2. Configure environment variables:
   ```sh
   export AUCTION_INTERVAL="5m" # Auction interval time
   export AUCTION_DURATION="1h" # Auction duration
   export MONGO_URI="mongodb://localhost:27017" # MongoDB URI
   export BATCH_INSERT_INTERVAL="20s" # Batch insert interval
   export MAX_BATCH_SIZE=4 # Maximum batch size
   export MONGO_INITDB_ROOT_USERNAME="admin" # MongoDB root username
   export MONGO_INITDB_ROOT_PASSWORD="admin" # MongoDB root password
   export MONGODB_URL="mongodb://admin:admin@mongodb:27017/auctions?authSource=admin" # MongoDB URL
   export MONGODB_DB="auctions" # MongoDB database name
   ```

3. Run the project:
   ```sh
   go run main.go
   ```

### Using Docker Compose

1. Ensure Docker is installed and running on your machine.

2. Run the following command to start the services:
   ```sh
   docker-compose up --build
   ```

This will build and start the application and MongoDB services as defined in the `docker-compose.yml` file.

### Pre-registered User ID for Testing

The following user is pre-registered in the MongoDB database for testing purposes:
- User ID: `a0d06633-4f0d-42ce-a653-d41a2d4aff94`
- User Name: `Admin`

## Environment Variables

- `AUCTION_INTERVAL`: Defines the auction interval time (e.g., "5m" for 5 minutes).
- `AUCTION_DURATION`: Defines the auction duration (e.g., "1h" for 1 hour).
- `MONGO_URI`: MongoDB connection URI (e.g., "mongodb://localhost:27017").
- `BATCH_INSERT_INTERVAL`: Interval for batch inserts (e.g., "20s" for 20 seconds).
- `MAX_BATCH_SIZE`: Maximum size for batch inserts.
- `MONGO_INITDB_ROOT_USERNAME`: MongoDB root username.
- `MONGO_INITDB_ROOT_PASSWORD`: MongoDB root password.
- `MONGODB_URL`: MongoDB connection URL (e.g., "mongodb://admin:admin@mongodb:27017/auctions?authSource=admin").
- `MONGODB_DB`: MongoDB database name.

## How to Run Tests

Run the tests:
   ```sh
   go test -v ./internal/usecase/auction_usecase/...
   ```

This will run all the tests in the project and display the results in the console.