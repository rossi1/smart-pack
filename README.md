# Smart-Pack Backend API

This is the backend service for **Smart-Pack** — an application that calculates optimal packaging solutions based on available pack sizes and the number of items ordered.

## What is Smart-Pack?

Smart-Pack solves the **bin packing optimization problem** by finding the minimum number of packages needed to fulfill an order. Similar to making change with coins, it determines the optimal combination of available pack sizes to minimize waste and total packages.

**Example**: If you need 263 items and have pack sizes of [250, 500, 1000, 2000, 5000]:
- Solution: 1×250 + 1×250 = 500 items total (2 packs, 237 excess items)
- This minimizes both the number of packs and excess items

## Algorithm

Smart-Pack uses a **dynamic programming approach** similar to the coin change problem:

1. **Minimize total packs first** - Find the smallest number of packages that can fulfill the order
2. **Minimize excess items second** - Among solutions with the same pack count, choose the one with least excess
3. **Configurable pack sizes** - Pack sizes can be updated without code changes

This ensures optimal packaging while allowing flexibility in pack size configuration.

## Features

* Calculate minimal total packs and items to fulfill an order
* Manage configurable pack sizes
* REST API endpoints for calculating and managing pack sizes
* Implements domain-driven design (DDD) and clean architecture principles
* Unit tested with mocking for easy maintenance

## API Endpoints

### Calculate Optimal Packing
```http
POST /api/v1/calculate
Content-Type: application/json

{
  "item_ordered": 263,
}
```

Response:

```json
{
  "total_packs": 2,
  "total_items": 500,
  "pack_breakdown": {
    "250": 2
  }
}
```

### Manage Pack Sizes
```http
POST /api/v1/pack-sizes
Content-Type: application/json

{
  "pack_sizes": [250, 500, 1000, 2000, 5000]
}
```

Response:

```json
{
  "message": "Pack sizes updated successfully"
}
```

### Get Pack Sizes
```http
GET /api/v1/pack-sizes
```

Response:

```json
{
  "pack_sizes": [250, 500, 1000, 2000, 5000]
}
```

### Health Check
```http
GET /api/v1/health
```

Response:

```json
{
  "status": "healthy"
}
```

### Ping
```http
GET /api/v1/ping
```

Response:

```json
{
  "message": "pong"
}
```

## Architecture Overview

The backend uses **Clean Architecture** / **Domain-Driven Design (DDD)** principles, structured into the following layers:

* **Domain**: Core business logic and entities (pack calculation algorithms)
* **Adapters**: Implementations of interfaces, including pack calculator logic
* **Ports**: Implementations of interfaces, including pack calculator logic
* **Application**: Application services (commands and queries)
* **App**: Implements **CQRS** pattern by separating **Commands** (write operations) and **Queries** (read operations). This layer orchestrates business logic and handles API requests and responses.
* **REST API**: HTTP server layer that handles API requests and responses
* **Mocks**: Automatically generated mocks for unit testing

This separation enforces single responsibility and enables easy testing and swapping of implementations.

## Getting Started

### Prerequisites

* Go 1.22+
* PostgreSQL 15+
* Docker & Docker Compose (optional, for containerized run)
* Make (optional, for simplified commands)

### App.env

It is mandatory to have the `app.env` file in your root folder.

### Configuration

`config/config.go` contains configuration variables loaded from `app.env` and environment variables using [viper](https://github.com/spf13/viper).

### Running Locally

Clone the repository and run the app:

```bash
git clone https://github.com/rossi1/smart-pack.git
cd smart-pack
make run
```

The API will be available at `http://localhost:8080`.

### Start Service

```bash
make run-docker-compose
```

### Database Migrations

Run database migrations **up** (apply all pending migrations):

```bash
go run main.go migrate up
```

Run database migrations **down** (rollback last migration):

```bash
go run main.go migrate down
```

### Using Docker

Build and run with Docker:

```bash
make docker-build
make docker-run
```

## Testing

```bash
make test
```

## Make Commands

* `make run` — Run the app locally using `go run main.go api`
* `make migrate-up` — Run database migrations up (`go run main.go migrate up`)
* `make migrate-down` — Rollback database migrations (`go run main.go migrate down`)
* `make test` — Run all tests with coverage report
* `make docker-build` — Build the Docker image
* `make docker-run` — Run the Docker container
* `make run-docker-compose` — Start services with Docker Compose

## License

[MIT](LICENSE)