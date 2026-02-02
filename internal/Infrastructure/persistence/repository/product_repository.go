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

func (r *PostgresProductRepository) UpdateInfo(ctx context.Context, p *product.Product) error {
	conn := tx.GetConn(ctx, r.db)

	_, err := conn.ExecContext(ctx, `
		UPDATE products
		SET name = $1, description = $2, status = $3, updated_at = $4
		WHERE id = $5
	`, p.Name(), p.Description(), p.Status(), p.UpdatedAt(), p.ID())

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (r *PostgresProductRepository) Delete(ctx context.Context, id int64) error {
	conn := tx.GetConn(ctx, r.db)

	_, err := conn.ExecContext(ctx, `
		DELETE FROM products WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (r *PostgresProductRepository) FindByID(ctx context.Context, id int64) (*product.Product, error) {
	conn := tx.GetConn(ctx, r.db)

	row := conn.QueryRowContext(ctx, `
		SELECT id, name, description, sku, status, stock_available, stock_reserved, created_at, updated_at
		FROM products WHERE id = $1
	`, id)

	var p product.Product

	// err := row.Scan(&p.id,
	// 	&p.name,
	// 	&p.description,
	// 	&p.sku,
	// 	&p.status,
	// 	&p.stock.available,
	// 	&p.stock.reserved,
	// 	p.createdAt,
	// 	p.updatedAt)

	// if err != nil {
	// 	return nil, fmt.Errorf("failed to find product by ID: %w", err)
	// }

	return &p, nil
}
