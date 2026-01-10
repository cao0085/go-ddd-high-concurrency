# Flash Sale Order System

A high-concurrency flash sale system built with Go, demonstrating solutions to common distributed system challenges.

## Core Challenges Addressed

1. **Overselling Prevention** - 100 items, 1000 buyers - guaranteed no overselling
2. **High Concurrency Writes** - Multiple goroutines safely decrementing stock
3. **Message Queue** - Kafka for request/processing decoupling and backpressure control
4. **Data Consistency** - Order, payment, and inventory consistency
5. **Monitoring** - Real-time metrics for goroutines, DB connections, throughput

## Tech Stack

- **Backend**: Go 1.23.4
- **Database**: PostgreSQL 17.2
- **Cache**: Redis 7.4.1
- **Message Queue**: Kafka 3.8.1 (KRaft mode)
- **Monitoring**: Prometheus 3.1.0 + Grafana 11.4.0

## Architecture

```
User Request → Go API
              ↓
         Redis (atomic stock decrement)
              ↓
         Kafka (order message)
              ↓
         Consumer → PostgreSQL (order persistence)
              ↓
         Periodic sync: Redis → PostgreSQL
```

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.23+ (for local development)
- Make (optional, for convenience commands)

### 1. Clone and Setup

```bash
cd flash-sale-order-system
cp .env.example .env
```

### 2. Start All Services

```bash
make up
```

Or without Make:
```bash
docker-compose up -d
```

### 3. Access Services

- **API**: http://localhost:8080
- **Grafana**: http://localhost:3000 (admin/admin123)
- **Prometheus**: http://localhost:9090
- **PostgreSQL**: localhost:5432 (flashsale/flashsale123)
- **Redis**: localhost:6379
- **Kafka**: localhost:9092

## Available Commands

```bash
make help        # Show all available commands
make up          # Start all services
make down        # Stop all services
make logs        # Show logs
make db-shell    # Connect to PostgreSQL
make redis-cli   # Connect to Redis
make clean       # Remove all containers and volumes
```

## Project Structure

```
.
├── cmd/
│   └── api/                 # Main application entry point
├── internal/
│   ├── handler/            # HTTP handlers
│   ├── service/            # Business logic
│   ├── repository/         # Database operations
│   ├── model/              # Data models
│   └── middleware/         # Middleware (logging, metrics, etc.)
├── pkg/
│   ├── database/           # PostgreSQL client
│   ├── redis/              # Redis client
│   ├── kafka/              # Kafka producer/consumer
│   └── metrics/            # Prometheus metrics
├── config/                  # Configuration files
├── scripts/                 # Database initialization scripts
├── monitoring/              # Prometheus & Grafana configs
├── docker-compose.yml
├── Dockerfile
└── Makefile
```

## Development

### Initialize Go Modules

```bash
make go-init
```

### Run in Development Mode

```bash
make dev
```

### Run Tests

```bash
make test
```

## Monitoring

Access Grafana at http://localhost:3000 to view:
- Request throughput
- Goroutine count
- Database connection pool
- Redis operations
- Kafka message rate
- Stock levels

## Key Implementation Details

### Preventing Overselling

1. **Redis Atomic Operations**: Use `DECR` for instant stock checks
2. **PostgreSQL Row Locks**: `SELECT FOR UPDATE` with transactions
3. **Optimistic Locking**: Version-based concurrency control

### High Concurrency

- Goroutine pool with max workers
- Connection pooling (PostgreSQL, Redis)
- Circuit breaker pattern for fault tolerance

### Data Consistency

- Two-phase commit for order + payment
- Idempotent message processing
- Periodic reconciliation (Redis ↔ PostgreSQL)

## License

MIT
