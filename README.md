# PricePulse 📈

A high-performance asynchronous price monitoring service built with Go. It tracks product prices, notifies about changes, and uses a microservices-ready architecture.

## 🚀 Features
- **Containerized Environment**: Fully Dockerized setup for seamless deployment.
- **Asynchronous Processing**: Uses Kafka to decouple API and price tracking logic.
- **Real-time Scraping**: Integrated with Colly for fetching actual prices from web pages.
- **Fast Caching**: Redis integration for quick price lookups and reduced DB load.
- **Clean Architecture**: Strictly separated layers: Domain, Service, Repository, and Transport.
- **Graceful Shutdown**: Properly closes all connections (Postgres, Kafka, Redis) to ensure data integrity.
- **Structured Logging**: JSON logs using `log/slog` for modern observability.
- **Static Analysis**: Integrated `golangci-lint` for maintaining high code quality.

## 🛠 Tech Stack
- **Language**: Go (Golang)
- **Containerization**: Docker & Docker Compose
- **Database**: PostgreSQL (Persistence)
- **Message Broker**: Apache Kafka (Async task distribution)
- **Cache**: Redis
- **Monitoring**: Prometheus (Metrics) & Grafana (Dashboards)
- **Quality**: golangci-lint
- **Frameworks**: Gin (HTTP), Colly (Web Scraping), pgx (Database Driver)

## 🏗 Architecture & Design
The project follows **Clean Architecture** and **Event-Driven Design**:
1. **API Service**: Accepts HTTP requests, persists product info, and produces events to Kafka.
2. **Watcher Service**: Consumes events, scrapes real prices, and updates the Database/Cache.
3. **Protobuf Contracts**: API definitions are described in `.proto` files for future gRPC implementation.

## 🚦 Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.21+

