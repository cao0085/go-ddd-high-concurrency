package provider

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"flash-sale-order-system/internal/Infrastructure/idgen"
	infraquery "flash-sale-order-system/internal/Infrastructure/persistence/query"
	persistence "flash-sale-order-system/internal/Infrastructure/persistence/repository"
	"flash-sale-order-system/internal/application/product/command"
	"flash-sale-order-system/internal/application/product/query"
	httpProduct "flash-sale-order-system/internal/interfaces/http/product"
)

type ProductHandlers struct {
	Command *httpProduct.CommandHandler
	Query   *httpProduct.QueryHandler
}

func NewProductHandlers(db *sql.DB, idGen *idgen.IDGenerator) *ProductHandlers {
	// Repositories (for Command side)
	productRepo := persistence.NewPostgresProductRepository(db)
	pricingRepo := persistence.NewPostgresProductPricingRepository(db)

	// Query Service (for Query side - no domain dependency)
	productQueryService := infraquery.NewPostgresProductQuery(db)

	// Command Handlers
	createHandler := command.NewCreateProductHandler(db, idGen, productRepo, pricingRepo)
	updateInfoHandler := command.NewUpdateProductInfoHandler(db, productRepo)
	removeHandler := command.NewRemoveProductHandler(db, productRepo)

	// Query Handlers
	getHandler := query.NewGetProductHandler(productQueryService)

	return &ProductHandlers{
		Command: httpProduct.NewCommandHandler(createHandler, updateInfoHandler, removeHandler),
		Query:   httpProduct.NewQueryHandler(getHandler),
	}
}

func RegisterProductRoutes(rg *gin.RouterGroup, handlers *ProductHandlers) {
	httpProduct.RegisterRoutes(rg, handlers.Command, handlers.Query)
}
