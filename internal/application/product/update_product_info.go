// internal/application/product/update_product_info.go
package product

type UpdateProductInfoCommand struct {
    ProductID int64
    Name      string
	Prices    map[string]float64 // currency code to amount
}

type UpdateProductInfoHandler struct {
    productRepo product.Repository
}

func NewUpdateProductInfoHandler(repo product.Repository) *UpdateProductInfoHandler {
	return &UpdateProductInfoHandler{productRepo: repo}
}

func (h *UpdateProductInfoHandler) Handle(ctx context.Context, cmd UpdateProductInfoCommand) error {
    // 1. 取得商品
    prod, err := h.productRepo.FindByID(ctx, cmd.ProductID)
    if err != nil {
        return err
    }

    prices := make(map[shared.Currency]shared.Money)
    for currencyStr, amount := range cmd.Prices {
        currency := shared.Currency(currencyStr)
        money, err := shared.NewMoney(amount, currency)
        if err != nil {
            return err
        }
        prices[currency] = money
    }
    
    // 2. 建立新的價格清單
    priceList, err := product.NewPriceList(cmd.Prices)
    if err != nil {
        return err
    }
    
    // 3. 更新商品資訊 (不影響庫存)
    if err := prod.UpdateInfo(cmd.Name, priceList); err != nil {
        return err
    }

	for currency, money := range cmd.Prices {
		prod.PriceList().SetPrice(currency, money)
	}
    
    // 4. 儲存 (Save 會儲存整個 Product)
    return h.productRepo.Save(ctx, prod)
}
