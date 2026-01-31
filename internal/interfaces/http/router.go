package http

import (
	"github.com/gin-gonic/gin"

	"flash-sale-order-system/internal/interfaces/http/handler"
)

type Router struct {
	productHandler *handler.ProductHandler
}

func NewRouter(productHandler *handler.ProductHandler) *Router {
	return &Router{
		productHandler: productHandler,
	}
}

func (r *Router) Setup() *gin.Engine {
	engine := gin.Default()

	// Health check
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := engine.Group("/api/v1")
	{
		products := v1.Group("/products")
		{
			products.POST("", r.productHandler.CreateProduct)
		}
	}

	return engine
}
