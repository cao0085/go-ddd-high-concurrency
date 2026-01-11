// internal/application/product/create_product.go
package product

import (
    "context"
    
    domain "flash-sale-order-system/internal/domain/product"
    shared "flash-sale-order-system/internal/shared/domain"
)

type CreateProductInfoCommand struct {
    Name           string
    InitialStock   int32
    Prices         map[string]float64 // currency code to amount
}

type CreateProductHandler struct {
    productRepo domain.Repository
}

func NewCreateProductHandler(repo domain.Repository) *CreateProductHandler {
    return &CreateProductHandler{productRepo: repo}
}

func (h *CreateProductHandler) Handle(ctx context.Context, cmd CreateProductInfoCommand) error {

    stock, err := domain.NewStock(cmd.InitialStock)
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
    
    priceList, err := domain.NewPriceList(prices)
    if err != nil {
        return err
    }
    
    prod, err := domain.NewProduct(cmd.Name, stock, priceList)
    if err != nil {
        return err
    }
    
    return h.productRepo.Save(ctx, prod)
}
