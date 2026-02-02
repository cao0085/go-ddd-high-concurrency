package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	tx "flash-sale-order-system/internal/Infrastructure/persistence/tx"
	product "flash-sale-order-system/internal/domain/product"
	shareddomain "flash-sale-order-system/internal/shared/domain"
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

func (r *PostgresProductPricingRepository) FindByProductID(ctx context.Context, productID int64) (*product.ProductPricing, error) {
	conn := tx.GetConn(ctx, r.db)

	rows, err := conn.QueryContext(ctx, `
		SELECT currency, amount, valid_from, valid_until
		FROM product_pricing
		WHERE product_id = $1
		ORDER BY valid_from, currency
	`, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query product pricing: %w", err)
	}
	defer rows.Close()

	// 用 (valid_from, valid_until) 作為 key 來分組
	type periodKey struct {
		validFrom  time.Time
		validUntil *time.Time
	}
	periodPrices := make(map[periodKey]map[shareddomain.Currency]shareddomain.Money)

	for rows.Next() {
		var (
			currency   string
			amount     float64
			validFrom  time.Time
			validUntil sql.NullTime
		)

		if err := rows.Scan(&currency, &amount, &validFrom, &validUntil); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var until *time.Time
		if validUntil.Valid {
			until = &validUntil.Time
		}

		key := periodKey{validFrom: validFrom, validUntil: until}
		if periodPrices[key] == nil {
			periodPrices[key] = make(map[shareddomain.Currency]shareddomain.Money)
		}

		money, err := shareddomain.NewMoney(amount, shareddomain.Currency(currency))
		if err != nil {
			return nil, fmt.Errorf("failed to create money: %w", err)
		}
		periodPrices[key][shareddomain.Currency(currency)] = money
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	if len(periodPrices) == 0 {
		return nil, nil
	}

	// 組裝 ProductPricing aggregate
	pricing := product.NewProductPricing(productID)
	for key, prices := range periodPrices {
		multiPrice, err := shareddomain.NewMultiCurrencyPrice(prices)
		if err != nil {
			return nil, fmt.Errorf("failed to create multi currency price: %w", err)
		}

		if err := pricing.AddPeriod(multiPrice, key.validFrom, key.validUntil); err != nil {
			return nil, fmt.Errorf("failed to add period: %w", err)
		}
	}

	return pricing, nil
}
