package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	product "flash-sale-order-system/internal/domain/product"
)

type PostgresProductRepository struct {
    db *sql.DB
}

// NewPostgresProductRepository creates a new PostgresProductRepository
func NewPostgresProductRepository(db *sql.DB) *PostgresProductRepository {
    return &PostgresProductRepository{db: db}
}

// Save 儲存整個 Product (包含所有欄位)
func (r *PostgresProductRepository) Save(ctx context.Context, p *product.Product) error {
    // 如果有 ID,就是更新;沒有 ID,就是新增
    if p.ID() == 0 {
        return r.insert(ctx, p)
    }
    return r.update(ctx, p)
}

func (r *PostgresProductRepository) insert(ctx context.Context, p *product.Product) error {
    // 開始交易
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    // 1. 插入商品主表
    var productID int64
    err = tx.QueryRowContext(ctx, `
        INSERT INTO products (name, stock_available, stock_reserved, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, p.Name(), p.Stock().Available(), p.Stock().Reserved(), p.CreatedAt(), p.UpdatedAt()).Scan(&productID)

    if err != nil {
        return fmt.Errorf("failed to insert product: %w", err)
    }

    // 2. 插入價格清單
    for currency, money := range p.PriceList().GetAllPrices() {
        _, err = tx.ExecContext(ctx, `
            INSERT INTO product_prices (product_id, currency, amount)
            VALUES ($1, $2, $3)
        `, productID, currency, money.Amount())

        if err != nil {
            return fmt.Errorf("failed to insert product price for currency %s: %w", currency, err)
        }
    }

    // 3. 提交交易
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}

func (r *PostgresProductRepository) update(ctx context.Context, p *product.Product) error {
    tx, _ := r.db.BeginTx(ctx, nil)
    defer tx.Rollback()
    
    // 1. 更新商品主表 (包含庫存!)
    _, err := tx.ExecContext(ctx, `
        UPDATE products 
        SET name = $1, 
            stock_available = $2, 
            stock_reserved = $3, 
            updated_at = $4
        WHERE id = $5
    `, p.Name(), p.Stock().Available(), p.Stock().Reserved(), p.UpdatedAt(), p.ID())
    
    if err != nil {
        return err
    }
    
    // 2. 更新價格 (先刪除舊的,再插入新的)
    _, err = tx.ExecContext(ctx, `DELETE FROM product_prices WHERE product_id = $1`, p.ID())
    if err != nil {
        return err
    }
    
    for currency, money := range p.PriceList().GetAllPrices() {
        _, err = tx.ExecContext(ctx, `
            INSERT INTO product_prices (product_id, currency, amount)
            VALUES ($1, $2, $3)
        `, p.ID(), currency, money.Amount())
        
        if err != nil {
            return err
        }
    }
    
    return tx.Commit()
}
