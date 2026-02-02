package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
		INSERT INTO products (id, sku, name, description, status, available_stock, reserved_stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, p.ID(), p.SKU(), p.Name(), p.Description(), p.Status(), p.Stock().Available(), p.Stock().Reserved(), p.CreatedAt(), p.UpdatedAt())

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
		SELECT id, sku, name, description, status, available_stock, reserved_stock, created_at, updated_at
		FROM products WHERE id = $1
	`, id)

	var (
		pID            int64
		sku            string
		name           string
		description    sql.NullString
		status         int8
		stockAvailable int32
		stockReserved  int32
		createdAt      time.Time
		updatedAt      time.Time
	)

	err := row.Scan(&pID, &sku, &name, &description, &status, &stockAvailable, &stockReserved, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find product by ID: %w", err)
	}

	return product.ReconstructProduct(
		pID,
		sku,
		name,
		description.String,
		status,
		stockAvailable,
		stockReserved,
		createdAt,
		updatedAt,
	), nil
}
