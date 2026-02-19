# PricePulse 📈

A high-performance asynchronous price monitoring service built with Go. 

## 🚀 Features
- **Containerized Environment**: Fully Dockerized setup.
- **Asynchronous Processing**: Uses Kafka for price tracking events.
- **Fast Caching**: Redis integration for quick price lookups.
- **Clean Architecture**: Separated Domain, Service, Repository, and Transport layers.
- **Graceful Shutdown**: Safe termination of all system components.
- **Structured Logging**: JSON logs using `log/slog`.
- **API Documentation**: Interactive documentation via **Swagger UI**.
- **Static Analysis**: Enforced code quality with `golangci-lint`.

## 🛠 Tech Stack
- **Language**: Go (Golang)
- **Infrastructure**: Docker, PostgreSQL, Kafka, Redis
- **Monitoring**: Prometheus
- **API**: Gin Gonic, Swagger (swaggo)
- **Quality**: golangci-lint

## 🚦 Getting Started

### API Documentation
Once the app is running, access the interactive Swagger UI at:
`http://localhost:8080/swagger/index.html`

### Installation
1. Setup your `.env` file (copy from `.env.example`).
2. Run infrastructure: `docker-compose up -d postgres redis kafka`.
3. Run the app: `go run cmd/main.go` or `make run`.
4. Lint code: `make lint`.

## 📡 Roadmap
- [ ] **Protobuf & gRPC**: Define service contracts and implement gRPC server.
- [ ] **Unit Testing**: Implement mocks and achieve 80% coverage.
- [ ] **Grafana**: Setup visual dashboards for Prometheus metrics.