package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/health", h.handle)
}

func (h *HealthHandler) handle(c *gin.Context) {
	response := gin.H{
		"success": true,
	}

	c.JSON(http.StatusOK, response)
}
