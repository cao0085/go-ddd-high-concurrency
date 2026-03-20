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

func (r *PostgresProductPricingRepository) FindByProductID(ctx context.Context, productID int64) (*product.ProductPricing, error) {
	conn := tx.GetConn(ctx, r.db)

	rows, err := conn.QueryContext(ctx, `
		SELECT currency, amount, valid_from, valid_until
		FROM product_pricing WHERE product_id = $1
		ORDER BY valid_from
	`, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query pricing for product %d: %w", productID, err)
	}
	defer rows.Close()

	pricing := product.NewProductPricing(productID)

	// group by period (valid_from, valid_until)
	type row struct {
		currency   string
		amount     float64
		validFrom  time.Time
		validUntil sql.NullTime
	}
	var allRows []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.currency, &r.amount, &r.validFrom, &r.validUntil); err != nil {
			return nil, fmt.Errorf("failed to scan pricing row: %w", err)
		}
		allRows = append(allRows, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// group by (validFrom, validUntil) to build MultiCurrencyPrice
	type periodKey struct {
		from  time.Time
		until sql.NullTime
	}
	grouped := make(map[periodKey]map[shareddomain.Currency]shareddomain.Money)
	var order []periodKey
	for _, r := range allRows {
		key := periodKey{from: r.validFrom, until: r.validUntil}
		if _, ok := grouped[key]; !ok {
			grouped[key] = make(map[shareddomain.Currency]shareddomain.Money)
			order = append(order, key)
		}
		currency := shareddomain.Currency(r.currency)
		money, err := shareddomain.NewMoney(r.amount, currency)
		if err != nil {
			return nil, err
		}
		grouped[key][currency] = money
	}

	for _, key := range order {
		prices, err := shareddomain.NewMultiCurrencyPrice(grouped[key])
		if err != nil {
			return nil, err
		}
		var until *time.Time
		if key.until.Valid {
			t := key.until.Time
			until = &t
		}
		if err := pricing.AddPeriod(prices, key.from, until); err != nil {
			return nil, err
		}
	}

	return pricing, nil
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
