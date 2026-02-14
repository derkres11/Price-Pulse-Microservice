APP_NAME = pricepulse
DOCKER_CONFIG=docker-compose.yml

.PHONY: run stop build clean logs

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