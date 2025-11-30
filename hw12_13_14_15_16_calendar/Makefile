DOCKER_IMG="calendar:develop"
PROJECT_DIR := $(CURDIR)
MIGRATIONS_DIR := $(PROJECT_DIR)/migrations
BUILDER := migrate-builder
POSTGRES := postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o ./bin/calendar -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o ./bin/calendar_scheduler -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o ./bin/calendar_sender -ldflags "$(LDFLAGS)" ./cmd/sender


run: build
	./bin/calendar -config ./configs/config.yaml
	./bin/calendar_scheduler -config ./configs/config.yaml
	./bin/calendar_sender -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	./bin/calendar version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.63.4

lint: install-lint-deps
	golangci-lint run ./...

create-migration:
	docker run --rm -v ./migrations:/app/migrations $(BUILDER) create -ext sql -dir /app/migrations $(name)

migrate:
	docker run --rm -v ./migrations:/app/migrations $(BUILDER) -database $(POSTGRES) -path /app/migrations up

generate:
	"C:\protoc\bin\protoc.exe" \
	--go_out=calendar/pb \
	--go-grpc_out=calendar/pb \
	api/*.proto

.PHONY: build run build-img run-img version test lint create-migration migrate generate