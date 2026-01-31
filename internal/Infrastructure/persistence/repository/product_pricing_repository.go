package persistence

import (
	"context"
	"database/sql"
	"fmt"

	product "flash-sale-order-system/internal/domain/product"
	"flash-sale-order-system/internal/Infrastructure/persistence/tx"
)

type PostgresProductPricingRepository struct {
	db *sql.DB
}

func NewPostgresProductPricingRepository(db *sql.DB) product.ProductPricingRepository {
	return &PostgresProductPricingRepository{db: db}
}

func (r *PostgresProductPricingRepository) Save(ctx context.Context, pricing *product.ProductPricing) error {
	conn := tx.GetConn(ctx, r.db)

	for _, period := range pricing.Periods() {
		// 插入 price_period
		_, err := conn.ExecContext(ctx, `
			INSERT INTO product_price_periods (product_id, price_from, price_until)
			VALUES ($1, $2, $3)
		`, pricing.ProductID(), period.From(), period.Until())

		if err != nil {
			return fmt.Errorf("failed to insert price period: %w", err)
		}

		// 插入每個幣別的價格
		for currency, money := range period.Prices().GetAllPrices() {
			_, err = conn.ExecContext(ctx, `
				INSERT INTO product_prices (product_id, currency, amount, price_from)
				VALUES ($1, $2, $3, $4)
			`, pricing.ProductID(), currency, money.Amount(), period.From())

			if err != nil {
				return fmt.Errorf("failed to insert price for currency %s: %w", currency, err)
			}
		}
	}

	return nil
}
