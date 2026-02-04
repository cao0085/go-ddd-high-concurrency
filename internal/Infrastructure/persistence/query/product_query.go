package query

import (
	"context"
	"database/sql"

	appquery "flash-sale-order-system/internal/application/product/query"
)

type PostgresProductQuery struct {
	db *sql.DB
}

func NewPostgresProductQuery(db *sql.DB) appquery.ProductQueryService {
	return &PostgresProductQuery{db: db}
}

func (q *PostgresProductQuery) GetByID(ctx context.Context, id int64) (*appquery.ProductDTO, error) {
	row := q.db.QueryRowContext(ctx, `
		SELECT id, name, description, sku, status, stock_available, stock_reserved, created_at, updated_at
		FROM products WHERE id = $1
	`, id)

	var dto appquery.ProductDTO
	err := row.Scan(
		&dto.ID,
		&dto.Name,
		&dto.Description,
		&dto.SKU,
		&dto.Status,
		&dto.Stock.Available,
		&dto.Stock.Reserved,
		&dto.CreatedAt,
		&dto.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}

func (q *PostgresProductQuery) GetWithCurrentPrice(ctx context.Context, id int64) (*appquery.ProductWithPriceDTO, error) {
	row := q.db.QueryRowContext(ctx, `
		SELECT
			p.id, p.name, p.description, p.sku, p.status,
			p.stock_available, p.stock_reserved, p.created_at, p.updated_at,
			pp.amount, pp.currency
		FROM products p
		LEFT JOIN product_prices pp ON p.id = pp.product_id
			AND pp.valid_from <= NOW()
			AND (pp.valid_until IS NULL OR pp.valid_until > NOW())
		WHERE p.id = $1
	`, id)

	var dto appquery.ProductWithPriceDTO
	var amount sql.NullFloat64
	var currency sql.NullString

	err := row.Scan(
		&dto.ID,
		&dto.Name,
		&dto.Description,
		&dto.SKU,
		&dto.Status,
		&dto.Stock.Available,
		&dto.Stock.Reserved,
		&dto.CreatedAt,
		&dto.UpdatedAt,
		&amount,
		&currency,
	)
	if err != nil {
		return nil, err
	}

	if amount.Valid && currency.Valid {
		dto.CurrentPrice = &appquery.PriceDTO{
			Amount:   amount.Float64,
			Currency: currency.String,
		}
	}

	return &dto, nil
}
