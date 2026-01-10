.PHONY: help build up down restart logs clean test

# Default target
help:
	@echo "Flash Sale Order System - Available Commands:"
	@echo ""
	@echo "  make build      - Build all Docker images"
	@echo "  make up         - Start all services"
	@echo "  make down       - Stop all services"
	@echo "  make restart    - Restart all services"
	@echo "  make logs       - Show logs from all services"
	@echo "  make clean      - Remove all containers and volumes"
	@echo "  make test       - Run tests"
	@echo "  make db-shell   - Connect to PostgreSQL shell"
	@echo "  make redis-cli  - Connect to Redis CLI"
	@echo "  make kafka-logs - Show Kafka logs"
	@echo ""

# Build Docker images
build:
	@echo "Building Docker images..."
	docker-compose build

# Start all services
up:
	@echo "Starting all services..."
	docker-compose up -d postgres redis kafka prometheus grafana
	@echo ""
	@echo "Infrastructure services started!"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis:      localhost:6379"
	@echo "  - Kafka:      localhost:9092"
	@echo "  - Prometheus: http://localhost:9090"
	@echo "  - Grafana:    http://localhost:3000 (admin/admin123)"
	@echo ""
	@echo "To build and start the app: make app-up"

# Start all services including app
up-all:
	@echo "Starting all services including app..."
	docker-compose up -d
	@echo ""
	@echo "All services started! Access:"
	@echo "  - API:        http://localhost:8080"
	@echo "  - Grafana:    http://localhost:3000 (admin/admin123)"
	@echo "  - Prometheus: http://localhost:9090"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis:      localhost:6379"
	@echo "  - Kafka:      localhost:9092"

# Build and start app only
app-up:
	@echo "Building and starting app..."
	docker-compose up -d --build app
	@echo "App started at http://localhost:8080"

# Stop all services
down:
	@echo "Stopping all services..."
	docker-compose down

# Restart all services
restart:
	@echo "Restarting all services..."
	docker-compose restart

# Show logs
logs:
	docker-compose logs -f

# Clean up everything
clean:
	@echo "Removing all containers and volumes..."
	docker-compose down -v
	@echo "Cleanup complete!"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# PostgreSQL shell
db-shell:
	docker exec -it flashsale-postgres psql -U flashsale -d flashsale_db

# Redis CLI
redis-cli:
	docker exec -it flashsale-redis redis-cli

# Kafka logs
kafka-logs:
	docker logs -f flashsale-kafka

# Check service health
health:
	@echo "Checking service health..."
	@docker ps --filter "name=flashsale-" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# Initialize Go modules
go-init:
	@echo "Initializing Go modules..."
	go mod init github.com/yourusername/flash-sale-order-system
	go mod tidy

# Development mode (rebuild and restart app only)
dev:
	docker-compose up -d --build app
	docker-compose logs -f app
