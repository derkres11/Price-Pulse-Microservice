package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/derkres11/price-pulse/internal/domain"
	"github.com/derkres11/price-pulse/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	services *service.ProductService
	logger   *slog.Logger
}

// Update constructor to accept logger
func NewHandler(services *service.ProductService, logger *slog.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	products := router.Group("/products")
	{
		products.POST("/", h.CreateProduct)
		products.GET("/:id", h.GetProduct)
	}

	return router
}

// CreateProduct handler
func (h *Handler) CreateProduct(c *gin.Context) {
	var input domain.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("invalid input", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.services.Create(c.Request.Context(), &input); err != nil {
		h.logger.Error("failed to create product", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("product created successfully", slog.String("url", input.URL))
	c.JSON(http.StatusCreated, input)
}

// GetProduct handler
func (h *Handler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	product, err := h.services.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get product", slog.Int64("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}
