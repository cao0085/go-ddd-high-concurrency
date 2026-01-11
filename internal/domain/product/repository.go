package product

import "context"

type Repository interface {
    // 基本 CRUD
    FindByID(ctx context.Context, id int64) (*Product, error)
    Save(ctx context.Context, p *Product) error  // 儲存整個 Product (包含所有欄位)
    Delete(ctx context.Context, id int64) error
    
    // 查詢方法
    FindByIDs(ctx context.Context, ids []int64) ([]*Product, error)
    ExistsByID(ctx context.Context, id int64) (bool, error)

    // 庫存操作
    UpdateStock(ctx context.Context, productID int64, newStock Stock) error
}