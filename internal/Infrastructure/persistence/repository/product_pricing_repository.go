package persistence

import (
	"context"
	"database/sql"
	"fmt"

	tx "flash-sale-order-system/internal/Infrastructure/persistence/tx"
	product "flash-sale-order-system/internal/domain/product"
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
		for currency, money := range period.Prices().GetAllPrices() {
			_, err := conn.ExecContext(ctx, `
				INSERT INTO product_pricing (product_id, currency, amount, valid_from, valid_until)
				VALUES ($1, $2, $3, $4, $5)
			`, pricing.ProductID(), currency, money.Amount(), period.ValidFrom(), period.ValidUntil())

			if err != nil {
				return fmt.Errorf("failed to insert price for currency %s: %w", currency, err)
			}
		}
	}

	return nil
}
