# Company Service

This is the gRPC-based `CompanyService` that provides APIs to manage company information, including creating, updating, deleting, and retrieving companies.

---

## Features

- **Create**: Add a new company with detailed information.
- **Patch**: Update specific fields of an existing company.
- **Delete**: Remove a company by its ID.
- **Get**: Retrieve a company by its ID.

---

## Prerequisites

1. **Install Go**: Ensure you have Go installed (version 1.20 or later).
   - [Download Go](https://golang.org/dl/)

2. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

3. **Protocol Buffers**:
   - Install the `protoc` compiler if you haven’t already.
   - Generate the Go code from `.proto` files:
     ```bash
     protoc --go_out=. --go-grpc_out=. proto/*.proto
     ```

4. **Docker**:
   - Ensure Docker is installed and running on your system.
   - [Download Docker](https://www.docker.com/products/docker-desktop/)

---

## Running the Main Application

   You can use the provided `docker-compose.yml` file to start all required services (Kafka, Zookeeper, Postgres, and the CompanyService):
   ```bash
   ./init.sh
   ```
   This script will:
   - Start all services defined in `docker-compose.yml`.
   - Create the Kafka topic `events`.

---

## Running Tests

**Run Unit and Integration Tests**:
   ```bash
   go test ./... -v
   ```

---

## Environment Variables

| Variable            | Default       | Description                     |
|---------------------|---------------|---------------------------------|
| `SERVER_ADDRESS`    | `8080`        | Port where the server runs.     |
| `JWT_SECRET`        | `your-secret` | Secret key for JWT authentication |
| `KAFKA_BROKERS`     | `kafka:9092`  | Kafka broker addresses          |
| `DESTINATION_TOPIC` | `events`      | Kafka topic for events          |
| `DB_HOST`           | `postgres`    | Database host                   |
| `DB_PORT`           | `5432`        | Database port                   |
| `DB_USER`           | `myuser`      | Database user                   |
| `DB_PASSWORD`       | `mypassword`  | Database password               |
| `DB_NAME`           | `mydb`        | Database name                   |

---

## Dependencies

- [GORM](https://gorm.io/) for database ORM.
- [Protocol Buffers](https://developers.google.com/protocol-buffers) for gRPC definitions.
- [golang-jwt](https://github.com/golang-jwt/jwt/v5) for JWT authentication.
- [grpcurl](https://github.com/fullstorydev/grpcurl) for testing gRPC services.
- [Confluent Kafka](https://www.confluent.io/) for message streaming.
- [PostgreSQL](https://www.postgresql.org/) as the database.

---

