package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (h *HealthController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/health", h.handle)
}

func (h *HealthController) handle(c *gin.Context) {
	response := gin.H{
		"success": true,
	}

	c.JSON(http.StatusOK, response)
}
