# Flash Sale Order System

A high-concurrency flash sale system built with Go, demonstrating solutions to common distributed system challenges.

## Core Challenges

1. **Overselling Prevention** - Lua script atomic decrement on Redis; no distributed lock needed
2. **High Concurrency Writes** - Redis handles stock at ~13,000 req/sec; DB writes are async via RabbitMQ
3. **Domain Decoupling** - Order placement only touches Redis; no cross-domain DB transaction required
4. **Data Consistency** - Warm-up syncs DB → Redis on startup; reconciliation goroutine syncs Redis → DB every 5 min
5. **Monitoring** - Real-time metrics via Prometheus + Grafana (app, Redis, RabbitMQ)

## Tech Stack

- **Backend**: Go 1.24.0
- **Database**: PostgreSQL 17.2
- **Cache**: Redis 7.4.1
- **Message Queue**: RabbitMQ 3.13 (with Management UI)
- **Monitoring**: Prometheus 3.1.0 + Grafana 11.4.0

## Features

- **DDD Architecture** - Domain / Application / Infrastructure / Interface layers; Product & Order aggregates, value objects (Money, Stock, Currency), Snowflake ID, repository pattern
- **Docker Infrastructure** - PostgreSQL, Redis, RabbitMQ, Prometheus, Grafana
- **Redis Stock Cache** - Atomic stock reservation via Lua script (available / reserved keys); overselling prevented without distributed lock
- **RabbitMQ Integration** - Order producer + consumer on **separate AMQP channels**, durable `orders` queue
- **Order Flow** - HTTP → Redis decrement → RabbitMQ publish → Consumer → PostgreSQL write → Redis ConfirmReservation
- **Redis Warm-up** - On startup, all product stock synced from PostgreSQL → Redis automatically
- **Periodic Reconciliation** - Background goroutine syncs Redis stock → PostgreSQL every 5 minutes

## Project Structure (Domain-Driven Design Architecture)

```md
.
├── cmd/
│   └── api/                          # Entry point, dependency injection
├── internal/
│   ├── domain/                       # Domain Layer
│   │   ├── product/                  # Product aggregate, value objects, repository interface
│   │   └── order/                    # Order aggregate (status lifecycle)
│   ├── application/                  # Application Layer
│   │   ├── product/
│   │   │   ├── command/              # CreateProduct, UpdateInfo, RemoveProduct
│   │   │   └── query/                # QueryProduct
│   │   ├── order/
│   │   │   ├── place_order.go        # Redis reserve → RabbitMQ publish
│   │   │   └── save_order.go         # Consumer: INSERT to DB → Redis confirm
│   │   └── stock/
│   │       └── service.go            # Warm-up & periodic reconciliation
│   ├── Infrastructure/               # Infrastructure Layer
│   │   ├── idgen/                    # Snowflake ID generator
│   │   └── persistence/
│   │       ├── postgres/             # DB connection
│   │       ├── redis/                # StockCache (Lua script atomic ops)
│   │       ├── repository/           # ProductRepository implementation
│   │       ├── query/                # Read-side product query
│   │       └── tx/                   # Transaction manager
│   ├── interfaces/                   # Interface Layer
│   │   └── http/
│   │       ├── order/                # Order handler & routes
│   │       ├── product/              # Product command/query handlers & routes
│   │       ├── middleware/           # Recovery middleware
│   │       └── router.go             # Gin router setup
│   ├── provider/                     # Wire product handlers
│   └── shared/domain/                # Shared value objects (Money, MultiCurrencyPrice)
├── pkg/
│   └── rabbitmq/                     # Connection, producer, consumer
├── scripts/
│   └── init.sql                      # Table creation + seed data
├── monitoring/                       # Prometheus & Grafana configs
├── docker-compose.yml
└── Dockerfile
```

## Flow

```md
Client
  │
  │  POST /api/v1/orders
  ▼
┌─────────────────────────────────┐
│   HTTP Handler (Gin)            │
│   PlaceOrder()                  │
│   • 驗證 JSON                   │
│   • 回傳 202 Accepted ◄──────── │ 
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
       │                                      成功 / 庫存不足
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
       │  │     成功 → Ack (移出 queue)  │
       │  │     失敗 → Nack + requeue   │
       │  └─────────────────────────────┘
       │
       └──③ Publish 失敗時 → CancelReservation()
                          Redis 庫存滾回
```

## Quick Start

### Prerequisites

- Docker & Docker Compose

### Start (Fresh)

```bash
docker compose down -v && docker compose up --build
```

> `-v` removes all volumes so PostgreSQL re-runs `init.sql` (creates tables + seeds sample data) and Redis starts clean.

### Access Services

- **API**: http://localhost:8080
- **RabbitMQ Management**: http://localhost:15672 (guest/guest)
- **Grafana**: http://localhost:3000 (admin/admin123)
- **Prometheus**: http://localhost:9090

## Verify Startup

### 1. Check DB — sample products seeded

```bash
docker exec flashsale-postgres psql -U flashsale -d flashsale_db \
  -c "SELECT id, name, available_stock, reserved_stock FROM products;"
```

```md
 id |           name            | available_stock | reserved_stock
----+---------------------------+-----------------+----------------
  1 | Flash Sale iPhone 15 Pro  |             100 |              0
  2 | Flash Sale MacBook Pro 16 |              50 |              0
  3 | Flash Sale AirPods Pro 2  |             200 |              0
```

### 2. Check Redis — warm-up synced from DB

```bash
docker exec flashsale-redis redis-cli MGET \
  stock:product:1:available stock:product:1:reserved \
  stock:product:2:available stock:product:2:reserved \
  stock:product:3:available stock:product:3:reserved
```

```md
1) "100"
2) "0"
3) "50"
4) "0"
5) "200"
6) "0"
```

### 3. Place Order

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "1",
    "product_id": 1,
    "quantity": 1,
    "price": 999.99
  }'
```

```json
{"order_id": 293227368682622976, "status": "queued"}
```

### 4. Verify Order Written to DB

```bash
docker exec flashsale-postgres psql -U flashsale -d flashsale_db \
  -c "SELECT id, order_id, user_id, product_id, quantity, status FROM orders;"
```

```md
 id |      order_id      | user_id | product_id | quantity | status
----+--------------------+---------+------------+----------+---------
  1 | 293227368682622976 |       1 |          1 |        1 | pending
```

---

## High Concurrency Testing

### 1. Install hey

```bash
brew install hey
```

### 2. Reset stock (Redis + DB)

```bash
docker exec flashsale-redis redis-cli SET "stock:product:1:available" 100000
docker exec flashsale-redis redis-cli SET "stock:product:1:reserved" 0
docker exec flashsale-postgres psql -U flashsale -d flashsale_db \
  -c "UPDATE products SET available_stock = 100000, reserved_stock = 0 WHERE id = 1;"
```

### 3. Run load test — 10,000 requests, 200 concurrent

```bash
hey -n 10000 -c 200 -m POST \
  -H "Content-Type: application/json" \
  -d '{"user_id":"1","product_id":1,"quantity":1,"price":99.99}' \
  http://localhost:8080/api/v1/orders
```

Maybe See:

```md
Summary:
  Requests/sec: ~13,000
  Average:      15ms
  P99:          57ms
  [202] 9951 responses
```

### 4. Verify no overselling

```bash
# Redis and DB should match
docker exec flashsale-redis redis-cli GET "stock:product:1:available"
docker exec flashsale-postgres psql -U flashsale -d flashsale_db \
  -c "SELECT available_stock, reserved_stock FROM products WHERE id = 1;"

# Order count should equal requests processed
docker exec flashsale-postgres psql -U flashsale -d flashsale_db \
  -c "SELECT COUNT(*) FROM orders;"
```
