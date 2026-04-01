APP_NAME=operation-admin-backend

.PHONY: tidy test build run

tidy:
	go mod tidy

test:
	go test ./...

build:
	go build -o bin/$(APP_NAME) ./cmd/server

run:
	go run ./cmd/server
