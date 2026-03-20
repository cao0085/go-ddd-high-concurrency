# System

- Backend: Go 1.24.0
- Message Queue: RabbitMQ 3.13 (with Management UI)
- Database: PostgreSQL 17.2
- Cache/Lock: Redis 7.4.1
- Monitoring: Prometheus 3.1.0 + Grafana 11.4.0
- Frontend: Command line

## Docker Containers (6 services)

1. **postgres:17.2-alpine** - Main database (inventory, orders, payments)
2. **redis:7.4.1-alpine** - Distributed cache + atomic operations
3. **rabbitmq:3.13-management-alpine** - Message queue (port 5672, management UI 15672)
4. **prom/prometheus:v3.1.0** - Metrics collection
5. **grafana/grafana:11.4.0** - Monitoring dashboard
6. **Go app** - Backend service (built from Dockerfile)

## rebuild

```bash
docker compose down && docker compose up --build
```

## PostgreSQL

```bash
docker exec flashsale-postgres psql -U flashsale -d flashsale_db
```

```bash
# 查看 tables
\dt

# 查看表結構
\d tableName

# 查看訂單
SELECT * FROM orders;
```

## Redis

init stock (啟動後需手動初始化，或由 warm-up 機制自動載入)

```bash
docker exec flashsale-redis redis-cli SET "stock:product:1:available" 100
docker exec flashsale-redis redis-cli SET "stock:product:1:reserved" 0
```

## RabbitMQ

- Management UI: http://localhost:15672 (guest/guest)
- Queue: `orders` (durable)

## curl

```bash
# health check
curl http://localhost:8080/health
```

```bash
# create product
curl -X POST http://localhost:8080/api/v1/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "description": "A test product",
    "sku": "TEST-001",
    "quantity": 50,
    "prices": {
      "USD": 99.99,
      "TWD": 3000,
      "JPY": 15000
    },
    "price_from": "2026-01-01T00:00:00Z"
  }'
```

```bash
# place order (flash sale)
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

## Order Flow

```
POST /api/v1/orders
  → Redis Lua script 原子扣庫存 (available - quantity, reserved + quantity)
  → RabbitMQ Publish (orders queue, persistent)
  → Consumer goroutine 消費訊息
  → INSERT INTO orders (PostgreSQL)
```

<!--
# practice
超賣問題 — 100 件商品，1000 人搶購，如何保證不超賣？（這是 PostgreSQL 事務和鎖的實戰）
高並發寫 — 多個 goroutine 同時扣庫存，如何避免死鎖和性能下降？
消息隊列 — 用 RabbitMQ 解耦請求和處理，學習背壓控制
一致性 — 訂單、支付、庫存的數據一致性
監控 — 觀察 goroutine 數量、數據庫連接、吞吐量 -->
