package http

import (
	"log/slog"
	"net/http"

	"github.com/derkres11/price-pulse/internal/domain"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service domain.ProductService
	logger  *slog.Logger
}

func NewHandler(s domain.ProductService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/health", h.HealthCheck)

	router.POST("/products", h.TrackProduct)

	return router
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) TrackProduct(c *gin.Context) {
	var input struct {
		URL         string  `json:"url"`
		TargetPrice float64 `json:"target_price"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.TrackProduct(c.Request.Context(), input.URL, input.TargetPrice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product tracked"})
}
