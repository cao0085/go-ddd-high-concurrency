# System

- Backend: Go 1.23.4 (2025 stable)
- Message Queue: Kafka 3.8.1 (KRaft mode, no Zookeeper needed)
- Database: PostgreSQL 17.2
- Cache/Lock: Redis 7.4.1
- Monitoring: Prometheus 3.1.0 + Grafana 11.4.0
- Frontend: Command line

## Docker Containers (6 services)

1. **postgres:17.2-alpine** - Main database (inventory, orders, payments)
2. **redis:7.4.1-alpine** - Distributed cache + atomic operations
3. **apache/kafka:3.8.1** - Message queue (KRaft mode, ARM64 compatible)
4. **prom/prometheus:v3.1.0** - Metrics collection
5. **grafana/grafana:11.4.0** - Monitoring dashboard
6. **Go app** - Backend service (built from Dockerfile)

## PostgreSQL

``` bash
psql -U flashsale -d postgres
```

``` bash
# 查看db
postgres=# \l

# 進入db
\c dbName
\dt

# 查看表結構
\d tableName
```



<!-- 
# practice
超賣問題 — 100 件商品，1000 人搶購，如何保證不超賣？（這是 PostgreSQL 事務和鎖的實戰）
高並發寫 — 多個 goroutine 同時扣庫存，如何避免死鎖和性能下降？
消息隊列 — 用 Kafka 解耦請求和處理，學習背壓控制
一致性 — 訂單、支付、庫存的數據一致性
監控 — 觀察 goroutine 數量、數據庫連接、吞吐量 -->
