package product

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"flash-sale-order-system/internal/application/product/command"
)

type CommandHandler struct {
	createHandler     *command.CreateProductHandler
	updateInfoHandler *command.UpdateProductInfoHandler
	removeHandler     *command.RemoveProductHandler
}

func NewCommandHandler(
	createHandler *command.CreateProductHandler,
	updateInfoHandler *command.UpdateProductInfoHandler,
	removeHandler *command.RemoveProductHandler,
) *CommandHandler {
	return &CommandHandler{
		createHandler:     createHandler,
		updateInfoHandler: updateInfoHandler,
		removeHandler:     removeHandler,
	}
}

func (h *CommandHandler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := command.CreateProductCommand{
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

func (h *CommandHandler) UpdateInfo(c *gin.Context) {
	var req UpdateProductInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := command.UpdateProductInfoCommand{
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

func (h *CommandHandler) Delete(c *gin.Context) {
	var req RemoveProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := command.RemoveProductCommand{
		Id: req.Id,
	}

	if err := h.removeHandler.Handle(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
