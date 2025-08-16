.PHONY: build build-cli build-termux build-web deploy dev

# ====================================================================================
# BUILD TARGETS
# ====================================================================================
build:
	@echo "Building all binaries..."
	@go build -v -o bin/cli cmd/cli/cli.go
	@go build -v -o bin/termux cmd/termux/termux.go
	@go build -v -o bin/web cmd/web/server.go
	@echo "Build complete."

build-cli:
	@echo "Building CLI..."
	@go build -v -o bin/cli cmd/cli/cli.go
	@echo "CLI build complete."

build-termux:
	@echo "Building Termux..."
	@go build -v -o bin/termux cmd/termux/termux.go
	@echo "Termux build complete."

build-web:
	@echo "Building web server..."
	@go build -v -o bin/web cmd/web/server.go
	@echo "Web server build complete."


# ====================================================================================
# DOCKER TARGETS
# ====================================================================================
deploy:
	docker compose down && docker compose up -d --build
dev:
	docker compose -f docker-compose.dev.yml down && docker compose -f docker-compose.dev.yml up --build
