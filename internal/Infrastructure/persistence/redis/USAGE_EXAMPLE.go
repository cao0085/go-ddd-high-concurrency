package redis

/*
使用範例：如何在你的應用中整合 Redis

## 1. 初始化 Redis Client (在 main.go 或 provider 中)

```go
import (
	redisInfra "your-project/internal/Infrastructure/persistence/redis"
)

func main() {
	// 初始化 Redis
	redisClient, err := redisInfra.NewClient(redisInfra.Config{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer redisInfra.CloseClient(redisClient)

	// 創建 Stock Cache 和 Lock
	stockCache := redisInfra.NewStockCache(redisClient)
	distLock := redisInfra.NewDistributedLock(redisClient)
}
```

## 2. 商品上架時同步庫存到 Redis

```go
func (s *ProductService) CreateProduct(ctx context.Context, product *Product) error {
	// 1. 儲存到資料庫
	if err := s.repo.Insert(ctx, product); err != nil {
		return err
	}

	// 2. 同步到 Redis
	stock := product.Stock()
	if err := s.stockCache.InitStock(ctx, product.ID(), stock.Available(), stock.Reserved()); err != nil {
		log.Printf("failed to init stock cache: %v", err)
		// 不影響主流程，只記錄錯誤
	}

	return nil
}
```

## 3. 搶購流程 - 使用 Redis 預扣庫存

```go
func (s *OrderService) CreateOrder(ctx context.Context, productID int64, quantity int32) error {
	// Step 1: 使用分散式鎖保護
	lockKey := fmt.Sprintf("product:%d", productID)

	err := s.distLock.WithLock(ctx, lockKey, 5*time.Second, func() error {
		// Step 2: Redis 原子性預扣庫存
		success, err := s.stockCache.Reserve(ctx, productID, quantity)
		if err != nil {
			return fmt.Errorf("reserve stock failed: %w", err)
		}

		if !success {
			return errors.New("insufficient stock")
		}

		// Step 3: 創建訂單（寫入資料庫）
		order := NewOrder(productID, quantity)
		if err := s.orderRepo.Insert(ctx, order); err != nil {
			// 失敗則歸還庫存
			s.stockCache.CancelReservation(ctx, productID, quantity)
			return err
		}

		return nil
	})

	return err
}
```

## 4. 訂單付款成功 - 確認預扣

```go
func (s *OrderService) ConfirmPayment(ctx context.Context, orderID int64) error {
	// 1. 查詢訂單
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 2. 確認 Redis 預扣 (reserved -> 0)
	if err := s.stockCache.ConfirmReservation(ctx, order.ProductID, order.Quantity); err != nil {
		return err
	}

	// 3. 更新資料庫庫存
	if err := s.productRepo.ConfirmReservation(ctx, order.ProductID, order.Quantity); err != nil {
		// 如果資料庫失敗，需要補償邏輯
		return err
	}

	return nil
}
```

## 5. 訂單取消/超時 - 歸還庫存

```go
func (s *OrderService) CancelOrder(ctx context.Context, orderID int64) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 1. Redis 歸還庫存 (reserved -> available)
	if err := s.stockCache.CancelReservation(ctx, order.ProductID, order.Quantity); err != nil {
		return err
	}

	// 2. 更新資料庫
	if err := s.productRepo.CancelReservation(ctx, order.ProductID, order.Quantity); err != nil {
		return err
	}

	return nil
}
```

## 6. 定期同步 Redis 與資料庫 (防止不一致)

```go
func (s *SyncService) SyncStockFromDB(ctx context.Context, productID int64) error {
	// 1. 從資料庫讀取最新庫存
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return err
	}

	// 2. 更新到 Redis
	stock := product.Stock()
	return s.stockCache.InitStock(ctx, productID, stock.Available(), stock.Reserved())
}
```

## 7. 使用分散式鎖的其他場景

```go
// 場景：防止重複提交訂單
func (s *OrderService) CreateOrderWithDedup(ctx context.Context, userID, productID int64) error {
	lockKey := fmt.Sprintf("order:user:%d:product:%d", userID, productID)

	// 嘗試獲取鎖，最多重試 3 次
	acquired, err := s.distLock.AcquireWithRetry(
		ctx,
		lockKey,
		10*time.Second,  // TTL
		3,               // maxRetries
		100*time.Millisecond, // retryInterval
	)

	if err != nil {
		return err
	}

	if !acquired {
		return errors.New("order already in progress")
	}

	defer s.distLock.Release(ctx, lockKey)

	// 執行訂單創建邏輯...
	return nil
}
```

## 重要注意事項

1. **Redis 是快取，不是唯一真相來源**
   - 資料庫才是最終的 source of truth
   - Redis 故障時要有降級方案

2. **定期同步**
   - 建議每隔一段時間從 DB 同步到 Redis
   - 或在 Redis miss 時從 DB 載入

3. **TTL 設定**
   - 避免 Redis 無限增長
   - 熱門商品可以設定較長 TTL

4. **錯誤處理**
   - Redis 操作失敗不應該讓整個流程失敗
   - 記錄日誌並降級到資料庫

5. **Lua Script 的優勢**
   - 保證原子性
   - 減少網路往返次數
   - 避免 race condition
*/

func Path(path string) string {

	var slices []string
	var tempStr string

	for _, ch := range path {
		if ch == '/' {
			//
		} else {
			tempStr += string(ch)
		}
	}
	if tempStr != "" {
		slices = append(slices, tempStr)
	}

}
