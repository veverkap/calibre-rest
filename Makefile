.PHONY: help build run test clean docker-build docker-run dev

help:
	@echo 'Usage:'
	@echo '  make build     - Build the Go binary'
	@echo '  make run       - Run the server in production mode'
	@echo '  make dev       - Run the server in development mode'
	@echo '  make test      - Run tests'
	@echo '  make clean     - Clean build artifacts'
	@echo '  make docker-build - Build Docker images'
	@echo '  make docker-run   - Run with docker-compose'

# Build the Go binary
build:
	go build -o calibre-rest .

# Run in production mode
run: build
	./calibre-rest

# Run in development mode
dev: build
	./calibre-rest --dev

# Run tests
test:
	go test -v

# Clean build artifacts
clean:
	rm -f calibre-rest
	go clean

# Build Docker images (Go version)
docker-build:
	docker build . -f docker/Dockerfile.go -t ghcr.io/kencx/calibre_rest:$(version)-go-app --target=app
	docker build . -f docker/Dockerfile.go -t ghcr.io/kencx/calibre_rest:$(version)-go-calibre --target=calibre

# Run with docker-compose
docker-run:
	docker compose up -d app

# Install Go dependencies
deps:
	go mod download
	go mod tidy

# Format Go code
fmt:
	go fmt ./...

# Lint Go code (requires golangci-lint)
lint:
	golangci-lint run

# Build for multiple platforms
build-cross:
	GOOS=linux GOARCH=amd64 go build -o calibre-rest-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o calibre-rest-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o calibre-rest-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o calibre-rest-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o calibre-rest-windows-amd64.exe .