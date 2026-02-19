package http

import (
	"log/slog"
	"net/http"
	"strconv"

	_ "github.com/derkres11/price-pulse/docs" // Import generated docs
	"github.com/derkres11/price-pulse/internal/domain"
	"github.com/derkres11/price-pulse/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.ProductService
	logger   *slog.Logger
}

func NewHandler(services *service.ProductService, logger *slog.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

// @title PricePulse API
// @version 1.0
// @description API Server for Price Monitoring Service
// @host localhost:8080
// @BasePath /

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

// CreateProduct godoc
// @Summary Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param input body domain.Product true "Product info"
// @Success 201 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Router /products [post]

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

// GetProduct godoc
// @Summary Get product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} domain.Product
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]

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
