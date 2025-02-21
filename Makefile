APP_NAME = loan-service
DOCKER_COMPOSE = docker-compose

.PHONY: all build up down logs ps test clean db-connect db-migrate

all: build up

build:
	$(DOCKER_COMPOSE) build

up:
	$(DOCKER_COMPOSE) up

down:
	$(DOCKER_COMPOSE) down

logs:
	$(DOCKER_COMPOSE) logs -f

ps:
	$(DOCKER_COMPOSE) ps

test:
	go test ./... -v

clean:
	rm -rf build/
	docker system prune -f

db-connect:
	docker-compose exec db psql -U loanuser -d loandb
