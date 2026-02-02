package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	productapp "flash-sale-order-system/internal/application/product"
)

type Handler struct {
	createHandler     *productapp.CreateProductHandler
	updateInfoHandler *productapp.UpdateProductInfoHandler
	removeHandler     *productapp.RemoveProductHandler
	getHandler        *productapp.GetProductHandler
}

func NewHandler(
	createHandler *productapp.CreateProductHandler,
	updateInfoHandler *productapp.UpdateProductInfoHandler,
	removeHandler *productapp.RemoveProductHandler,
	getHandler *productapp.GetProductHandler) *Handler {
	return &Handler{
		createHandler:     createHandler,
		updateInfoHandler: updateInfoHandler,
		removeHandler:     removeHandler,
		getHandler:        getHandler,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	products := rg.Group("/product")
	{
		products.POST("", h.Create)
		products.GET("/:id", h.GetByID)
		products.PUT("/:id", h.UpdateInfo)
		products.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := productapp.CreateProductCommand{
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Quantity:    req.Quantity,
		Prices:      req.Prices,
		PriceFrom:   req.PriceFrom,
		PriceUntil:  req.PriceUntil,
	}

	productID, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateProductResponse{ID: productID})
}

func (h *Handler) GetByID(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// 呼叫 handler
	product, err := h.getHandler.Handle(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *Handler) UpdateInfo(c *gin.Context) {
	var req UpdateProductInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := productapp.UpdateProductInfoCommand{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}

	err := h.updateInfoHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) Delete(c *gin.Context) {
	var req RemoveProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := productapp.RemoveProductCommand{
		Id: req.Id,
	}

	if err := h.removeHandler.Handle(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
