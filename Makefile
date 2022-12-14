DOCKER_COMPOSE ?= docker-compose

default_target: build

build:
	@${DOCKER_COMPOSE} build

start:
	@${DOCKER_COMPOSE} up -d

stop:
	@${DOCKER_COMPOSE} down

.PHONY: test
test:
	@${DOCKER_COMPOSE} -f docker-compose.test.yml build
	@${DOCKER_COMPOSE} -f docker-compose.test.yml up -d audioapi_test
	@${DOCKER_COMPOSE} -f docker-compose.test.yml up test-client
	@${DOCKER_COMPOSE} -f docker-compose.test.yml down -v