package persistence

import (
	"context"
	"database/sql"
	"fmt"

	tx "flash-sale-order-system/internal/Infrastructure/persistence/tx"
	product "flash-sale-order-system/internal/domain/product"
)

type PostgresProductRepository struct {
	db *sql.DB
}

// NewPostgresProductRepository creates a new PostgresProductRepository
func NewPostgresProductRepository(db *sql.DB) product.ProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) Insert(ctx context.Context, p *product.Product) error {
	conn := tx.GetConn(ctx, r.db)

	_, err := conn.ExecContext(ctx, `
		INSERT INTO products (id, name, description, sku, status, stock_available, stock_reserved, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, p.ID(), p.Name(), p.Description(), p.SKU(), p.Status(), p.Stock().Available(), p.Stock().Reserved(), p.CreatedAt(), p.UpdatedAt())

	if err != nil {
		return fmt.Errorf("failed to insert product: %w", err)
	}

	return nil
}
