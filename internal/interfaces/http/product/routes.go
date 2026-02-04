package product

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, cmd *CommandHandler, qry *QueryHandler) {
	products := rg.Group("/product")
	{
		// Query endpoints
		products.GET("/:id", qry.GetByID)

		// Command endpoints
		products.POST("", cmd.Create)
		products.PUT("/:id", cmd.UpdateInfo)
		products.DELETE("/:id", cmd.Delete)
	}
}
