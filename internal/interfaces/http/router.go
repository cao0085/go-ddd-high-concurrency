package http

import (
	"flash-sale-order-system/internal/interfaces/http/middleware"
	"flash-sale-order-system/internal/interfaces/http/product"

	"github.com/gin-gonic/gin"
)

type Router struct {
	handlers *Handlers
}

func NewRouter(handlers *Handlers) *Router {
	return &Router{handlers: handlers}
}

func (r *Router) Setup() *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.Recovery())

	// Health check
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := engine.Group("/api/v1")
	{
		product.RegisterRoutes(v1, r.handlers.ProductCommand, r.handlers.ProductQuery)
	}

	return engine
}
