.PHONY: all build run test clean docker-up docker-down

# ==============================================================================
# Backend Commands (Go)
# ==============================================================================

build-backend:
	cd backend && go mod tidy && go build -o bin/server cmd/server/main.go

run-backend:
	cd backend && go run cmd/server/main.go

test-backend:
	cd backend && go test -v -cover ./...

# ==============================================================================
# Frontend Commands (Node/npm)
# ==============================================================================

install-frontend:
	cd frontend && npm install

run-frontend:
	cd frontend && npm run dev

build-frontend:
	cd frontend && npm run build

# ==============================================================================
# Docker Infrastructure Commands
# ==============================================================================

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# ==============================================================================
# Init / Bootstrap
# ==============================================================================
init: docker-up
	@echo "Waiting for MySQL to initialize..."
	@sleep 10
	@echo "Infrastructure is up! Remember to pull ollama model: docker exec -it hztour-ollama ollama run qwen:4b"
