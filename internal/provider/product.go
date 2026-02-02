package provider

import (
	"database/sql"

	"flash-sale-order-system/internal/Infrastructure/idgen"
	persistence "flash-sale-order-system/internal/Infrastructure/persistence/repository"
	productapp "flash-sale-order-system/internal/application/product"
	httpProduct "flash-sale-order-system/internal/interfaces/http/product"
)

func NewProductHandler(db *sql.DB, idGen *idgen.IDGenerator) *httpProduct.Handler {
	productRepo := persistence.NewPostgresProductRepository(db)
	pricingRepo := persistence.NewPostgresProductPricingRepository(db)

	createHandler := productapp.NewCreateProductHandler(db, idGen, productRepo, pricingRepo)
	updateInfoHandler := productapp.NewUpdateProductInfoHandler(db, productRepo)
	getHandler := productapp.NewGetProductHandler(db, productRepo)
	removeHandler := productapp.NewRemoveProductHandler(db, productRepo)

	return httpProduct.NewHandler(createHandler, updateInfoHandler, removeHandler, getHandler)
}
