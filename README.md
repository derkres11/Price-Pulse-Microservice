# PricePulse 📈

A high-performance asynchronous price monitoring service built with Go. PricePulse is designed for scalability and low latency, utilizing a hybrid communication model and a robust event-driven architecture.

## 🚀 Key Features

* **Hybrid API Interface**: Support for both **REST (JSON)** and **gRPC** (Protobuf) for efficient inter-service communication.
* **Event-Driven Architecture**: Asynchronous processing using **Apache Kafka** for price tracking and updates.
* **High-Performance Caching**: Multi-level caching with **Redis** to minimize database load.
* **Clean Architecture**: Strict separation of concerns (Domain, Service, Transport layers) following SOLID principles.
* **Production-Ready**: Implementation of **Graceful Shutdown**, structured JSON logging (`slog`), and health checks.
* **Containerized**: Fully orchestrated environment using Docker and Docker Compose.

## 🛠 Tech Stack

* **Language**: Go (Golang)
* **Communication**: gRPC (Protobuf), REST (Gin Gonic)
* **Infrastructure**: PostgreSQL, Apache Kafka, Redis, Docker
* **Observability**: Prometheus, Structured Logging (`slog`)
* **Documentation**: Swagger UI (swaggo)
* **Quality Assurance**: golangci-lint, gofumpt

## 🏗 System Architecture

The project is structured to be easily maintainable and testable:

* `cmd/api/` — Application entry point and dependency injection.
* `internal/domain/` — Core business entities and interface definitions.
* `internal/service/` — Implementation of business logic (Use Cases).
* `internal/transport/` — Delivery layer (gRPC handlers and HTTP controllers).
* `pkg/api/` — Auto-generated gRPC stubs and Protobuf code.
* `proto/` — Service contracts defined in Protocol Buffers.

## 🚦 Getting Started

### API Documentation

Once the application is running, access the interactive Swagger UI to explore the REST endpoints:
`http://localhost:8080/swagger/index.html`

### Installation & Setup

1. **Clone and Prepare**:
```bash
git clone https://github.com/derkres11/price-pulse.git
cp .env.example .env

```


2. **Launch Infrastructure**:
```bash
docker-compose up -d

```


3. **Run the Service**:
```bash
go run cmd/api/main.go

```



## 📡 Roadmap & Future Improvements

* [ ] **Notification Engine**: Integration with Telegram/Email alerts for price hits.
* [ ] **Comprehensive Testing**: Implementing unit and integration tests with **Testify** and **Mockery**.
* [ ] **Observability**: Setting up **Grafana** dashboards to visualize Prometheus metrics.
* [ ] **CI/CD**: Automated deployment pipelines using GitHub Actions.
