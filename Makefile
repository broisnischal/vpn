.PHONY: build-server build-client build run-server run-client docker-build docker-up docker-down clean test

# Build server binary
build-server:
	go build -o bin/omail-server ./cmd/server

# Build client binary
build-client:
	go build -o bin/omail-client ./cmd/client

# Build both
build: build-server build-client

# Run server (requires root for TUN)
run-server:
	sudo ./bin/omail-server -address :51820 -password changeme123

# Run client (requires root for TUN)
run-client:
	sudo ./bin/omail-client -server localhost:51820 -password changeme123

# Docker operations
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run ./...
