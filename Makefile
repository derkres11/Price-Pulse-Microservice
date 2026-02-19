APP_NAME = pricepulse
DOCKER_CONFIG=docker-compose.yml
include .env
export

.PHONY: run stop build clean logs migrate-create migrate-up lint

run:
	docker-compose -f $(DOCKER_CONFIG) up -d

stop:
	docker-compose -f $(DOCKER_CONFIG) down

build:
	go build -o bin/$(APP_NAME) cmd/app/main.go

logs:
	docker-compose -f $(DOCKER_CONFIG) logs -f

clean:
	rm -rf bin/
	docker system prune -f

DB_URL=postgres://admin:secret@localhost:5432/pricepulse?sslmode=disable

migrate-create:
	docker run --rm -v $(shell pwd)/migrations:/migrations migrate/migrate create -ext sql -dir /migrations/ -seq $(name)

migrate-up:
	docker run --rm -v $(shell pwd)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "$(DB_URL)" up

lint:
	golangci-lint run ./...