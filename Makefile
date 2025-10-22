.PHONY: all build test lint audit integration-test env-up env-start env-stop env-down coverage clean guard-% destroy

NAME=gofxtest
MODULE=github.com/m2tx/${NAME}
TAG=$(shell git rev-parse --short HEAD)

ifneq (${GITHUB_REF_NAME},)
	TAG = ${GITHUB_REF_NAME}
endif

DOCKER_FILE = ./build/Dockerfile
DOCKER_NETWORK = ${NAME}_default
COMPOSE_FILE = ./build/docker-compose.yml

__docker-build: guard-DOCKER_BUILD_TARGET
	docker build --tag=$(NAME) --progress=plain --target=$(DOCKER_BUILD_TARGET) $(DOCKER_BUILD_ARGS) -f $(DOCKER_FILE) .

__docker-build-tools: guard-APP guard-APP_VERSION
	docker build --tag=$(NAME):$(APP)-$(APP_VERSION)  --build-arg VERSION=$(APP_VERSION) -f $(DOCKER_FILE).$(APP) .
	docker run --rm -v $(shell pwd):/app -w /app $(NAME):$(APP)-$(APP_VERSION)

__compose:
	DOCKER_BUILDKIT=1 \
	NETWORK_NAME=$(DOCKER_NETWORK) \
	docker compose --file $(COMPOSE_FILE) --project-name $(NAME) $(COMPOSE_COMMAND)

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Variable '$*' not set"; \
		exit 1; \
	fi 

env-up:
	@make __compose COMPOSE_COMMAND="up -d"

env-up-%:
	@make __compose COMPOSE_COMMAND="up -d $*"

env-start: 
	@make __compose COMPOSE_COMMAND="start"

env-stop: 
	@make __compose COMPOSE_COMMAND="stop"

env-down:
	@make __compose COMPOSE_COMMAND="down"

install:
	go mod vendor

install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest

swag:
	swag init -g ./internal/http/* -o ./docs

lint:
	APP=golangci \
	APP_VERSION=v2.5.0 \
	make __docker-build-tools

audit: audit-osvscanner

audit-gosec:
	APP=gosec \
	APP_VERSION=v2.22.9 \
	make __docker-build-tools

audit-osvscanner:
	APP=osvscanner \
	APP_VERSION=v2.2.2 \
	make __docker-build-tools

test: env-up-mongo test-run test-clean

test-run:
	DOCKER_BUILD_ARGS=" --build-arg NAME=${NAME} --build-arg MODULE=${MODULE} " \
	DOCKER_BUILD_TARGET="test" \
	make __docker-build
	mkdir -p coverage
	docker run -e MONGO_DATABASE="test" -e MONGO_URL="mongodb://mongo:27017" --network $(DOCKER_NETWORK) --name=$(NAME)-test $(NAME)
	docker cp $(NAME)-test:/coverage.html $(shell pwd)/coverage/coverage.html
	docker cp $(NAME)-test:/coverage.txt $(shell pwd)/coverage/coverage.txt
	docker cp $(NAME)-test:/coverage.out $(shell pwd)/coverage/coverage.out

test-clean:
	docker rm -fv $(NAME)-test

coverage: test
	sed -i 's/black/whitesmoke/g' $(shell pwd)/coverage/coverage.html
	open $(shell pwd)/coverage/coverage.html

build:
	DOCKER_BUILDKIT=1 \
	DOCKER_BUILD_ARGS=" --build-arg NAME=${NAME} --build-arg MODULE=${MODULE} " \
	DOCKER_BUILD_TARGET=build \
	make __docker-build

image:
	DOCKER_BUILDKIT=1 \
	DOCKER_BUILD_ARGS=" --build-arg NAME=${NAME} --build-arg MODULE=${MODULE} " \
	DOCKER_BUILD_TARGET=image \
	make __docker-build

publish: image
	docker login --username=${DOCKER_USERNAME} --password=${DOCKER_PASSWORD}
	docker image tag ${NAME} ${DOCKER_USERNAME}/${NAME}:${TAG}
	docker push ${DOCKER_USERNAME}/${NAME}:${TAG}

destroy: test-clean env-down

run: env-up
	export $(shell grep -v '^#' .env | xargs) && go run ./cmd/server/main.go	