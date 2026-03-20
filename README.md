# Flash Sale Order System (OnGoing)

A high-concurrency flash sale system built with Go, demonstrating solutions to common distributed system challenges.

## Progress

### ✅ Completed
- **DDD Domain Layer** - Product aggregate, value objects (Money, Stock, Currency), domain errors
- **DDD Application Layer** - Product use cases (create/update/delete) + Order placement use case
- **DDD Infrastructure Layer** - PostgreSQL repositories, transaction manager, Snowflake ID generator
- **DDD Interface Layer** - Gin HTTP router, handlers with request/response DTOs
- **Docker Infrastructure** - PostgreSQL, Redis, RabbitMQ, Prometheus, Grafana
- **Redis Stock Cache** - Atomic stock reservation via Lua script (available / reserved keys)
- **RabbitMQ Integration** - Order producer + consumer, durable `orders` queue
- **Order Flow** - HTTP → Redis decrement → RabbitMQ publish → Consumer → PostgreSQL write

### 🚧 Next Steps
- [ ] Implement full domain logic (User aggregate, Payment aggregate)
- [ ] Redis warm-up on startup (sync stock from PostgreSQL to Redis)
- [ ] Distributed lock for overselling prevention
- [ ] Idempotent order processing (dedup by order_id)
- [ ] Periodic reconciliation: Redis ↔ PostgreSQL

---

## Core Challenges Addressed

1. **Overselling Prevention** - 100 items, 1000 buyers - guaranteed no overselling
2. **High Concurrency Writes** - Atomic Lua script in Redis for race-free stock decrement
3. **Message Queue** - RabbitMQ for request/processing decoupling and backpressure control
4. **Data Consistency** - Order and inventory consistency via queue + consumer
5. **Monitoring** - Real-time metrics for goroutines, DB connections, throughput

## Tech Stack

- **Backend**: Go 1.24.0
- **Database**: PostgreSQL 17.2
- **Cache**: Redis 7.4.1
- **Message Queue**: RabbitMQ 3.13 (with Management UI)
- **Monitoring**: Prometheus 3.1.0 + Grafana 11.4.0

## Architecture

```md
Client
  │
  │  POST /api/v1/orders
  ▼
┌─────────────────────────────────┐
│   HTTP Handler (Gin)            │
│   PlaceOrder()                  │
│   • 驗證 JSON                   │
│   • 回傳 202 Accepted ◄──────── │ ← 這裡就結束了，不等 DB
└──────────────┬──────────────────┘
               │
               ▼
┌─────────────────────────────────┐
│   Application Layer             │
│   PlaceOrderHandler.Handle()    │
└──────┬──────────────────────────┘
       │
       ├──① Reserve()──────────────────────────────────┐
       │                                               ▼
       │                               ┌───────────────────────────┐
       │                               │  Redis (Lua Script 原子)   │
       │                               │  stock:product:1:available │
       │                               │  DECRBY 1                  │
       │                               │  stock:product:1:reserved  │
       │                               │  INCRBY 1                  │
       │                               └───────────┬───────────────┘
       │                                           │
       │                              ✅ 成功 / ❌ 庫存不足
       │
       ├──② Publish()─────────────────────────────────┐
       │   (扣庫存成功才執行)                          ▼
       │                               ┌───────────────────────────┐
       │                               │  RabbitMQ                  │
       │                               │  Queue: "orders"           │
       │                               │  Durable: true (持久化)    │
       │                               └───────────┬───────────────┘
       │                                           │
       │             ┌─────────────────────────────┘
       │             │  Consumer goroutine 監聽
       │             ▼
       │  ┌─────────────────────────────┐
       │  │  SaveOrderHandler.Handle()  │
       │  │  INSERT INTO orders         │
       │  │  status = 'pending'         │
       │  └─────────────┬───────────────┘
       │                │
       │                ▼
       │  ┌─────────────────────────────┐
       │  │  PostgreSQL                 │
       │  │  Table: orders              │
       │  │  ✅ 成功 → Ack (移出 queue)  │
       │  │  ❌ 失敗 → Nack + requeue   │
       │  └─────────────────────────────┘
       │
       └──③ Publish 失敗時 → CancelReservation()
                          Redis 庫存滾回
```

```
User Request → Go API
              ↓
         Redis (atomic stock decrement via Lua)
              ↓
         RabbitMQ (orders queue, durable)
              ↓
         Consumer goroutine → PostgreSQL (order persistence)
```

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.23+ (for local development)

### 1. Clone and Setup

```bash
cd flash-sale-order-system
cp .env.example .env
```

### 2. Start All Services

```bash
docker compose up --build
```

### 3. Initialize Redis Stock (required before placing orders)

```bash
docker exec flashsale-redis redis-cli SET "stock:product:1:available" 100
docker exec flashsale-redis redis-cli SET "stock:product:1:reserved" 0
```

### 4. Access Services

- **API**: http://localhost:8080
- **Grafana**: http://localhost:3000 (admin/admin123)
- **Prometheus**: http://localhost:9090
- **RabbitMQ Management**: http://localhost:15672 (guest/guest)
- **PostgreSQL**: localhost:5432 (flashsale/flashsale123)
- **Redis**: localhost:6379

## API Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```

### Product
```bash
# Create product
curl -X POST http://localhost:8080/api/v1/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "description": "A test product",
    "sku": "TEST-001",
    "quantity": 50,
    "prices": {"USD": 99.99, "TWD": 3000, "JPY": 15000},
    "price_from": "2026-01-01T00:00:00Z"
  }'

# Get product
curl http://localhost:8080/api/v1/product/1
```

### Order (Flash Sale)
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "ORD-001",
    "user_id": "1",
    "product_id": 1,
    "quantity": 1,
    "price": 999.99
  }'
```

Response: `{"order_id":"ORD-001","status":"queued"}`

### Verify Order in DB
```bash
docker exec flashsale-postgres psql -U flashsale -d flashsale_db -c "SELECT * FROM orders;"
```

## Project Structure (DDD)

```
.
├── cmd/
│   └── api/                          # Application entry point & DI
├── internal/
│   ├── domain/                       # Domain Layer
│   │   └── product/                  # Product aggregate
│   ├── application/                  # Application Layer
│   │   ├── product/                  # Product use cases
│   │   └── order/                    # Order use cases
│   │       ├── place_order.go        # Redis reserve + RabbitMQ publish
│   │       └── save_order.go         # Consumer: persist to PostgreSQL
│   ├── Infrastructure/               # Infrastructure Layer
│   │   ├── persistence/
│   │   │   ├── postgres/             # Database connection
│   │   │   ├── redis/                # Stock cache + distributed lock
│   │   │   ├── repository/           # Repository implementations
│   │   │   └── tx/                   # Transaction manager
│   │   └── idgen/                    # Snowflake ID generator
│   ├── interfaces/                   # Interface Layer
│   │   └── http/
│   │       ├── order/                # Order HTTP handler
│   │       ├── product/              # Product HTTP handler
│   │       ├── middleware/           # Middleware
│   │       └── router.go             # Gin router
│   └── shared/                       # Shared kernel
│       └── domain/                   # Shared value objects
├── pkg/
│   └── rabbitmq/                     # RabbitMQ connection, producer, consumer
├── scripts/                          # Database init scripts
├── monitoring/                       # Prometheus & Grafana configs
├── docker-compose.yml
└── Dockerfile
```

## License

MIT
