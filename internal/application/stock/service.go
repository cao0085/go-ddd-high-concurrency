package stock

import (
	"context"
	"log"
	"time"

	redisinfra "flash-sale-order-system/internal/Infrastructure/persistence/redis"
	"flash-sale-order-system/internal/domain/product"
)

// StockService handles Redis stock warm-up and periodic reconciliation.
type StockService struct {
	productRepo product.ProductRepository
	stockCache  *redisinfra.StockCache
}

// NewStockService creates a new StockService.
func NewStockService(productRepo product.ProductRepository, stockCache *redisinfra.StockCache) *StockService {
	return &StockService{
		productRepo: productRepo,
		stockCache:  stockCache,
	}
}

// WarmUp loads all product stock from PostgreSQL into Redis on startup.
func (s *StockService) WarmUp(ctx context.Context) error {
	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		return err
	}

	for _, p := range products {
		if err := s.stockCache.InitStock(ctx, p.ID(), p.Stock().Available(), p.Stock().Reserved()); err != nil {
			log.Printf("stock warm-up: failed to init product %d: %v", p.ID(), err)
			continue
		}
	}

	log.Printf("stock warm-up: initialized %d products", len(products))
	return nil
}

// StartReconcile starts a background goroutine that periodically syncs Redis stock back to PostgreSQL.
// Redis is treated as the source of truth during a flash sale.
// The goroutine exits when ctx is cancelled.
func (s *StockService) StartReconcile(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("stock reconcile: stopped")
			return
		case <-ticker.C:
			s.reconcile(ctx)
		}
	}
}

func (s *StockService) reconcile(ctx context.Context) {
	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		log.Printf("stock reconcile: failed to load products: %v", err)
		return
	}

	updated := 0
	for _, p := range products {
		redisAvail, redisReserved, err := s.stockCache.GetStock(ctx, p.ID())
		if err != nil {
			log.Printf("stock reconcile: failed to read redis for product %d: %v", p.ID(), err)
			continue
		}

		// Skip if both keys are missing (cache miss — product may not be active in flash sale)
		if redisAvail == 0 && redisReserved == 0 {
			continue
		}

		// Only update if Redis differs from DB
		if redisAvail == p.Stock().Available() && redisReserved == p.Stock().Reserved() {
			continue
		}

		if err := s.productRepo.UpdateStock(ctx, p.ID(), redisAvail, redisReserved); err != nil {
			log.Printf("stock reconcile: failed to update product %d: %v", p.ID(), err)
			continue
		}
		updated++
	}

	if updated > 0 {
		log.Printf("stock reconcile: synced %d products", updated)
	}
}
