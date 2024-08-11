package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DummyHandler struct {
}

func NewDummyHandler() *DummyHandler {
	return &DummyHandler{}
}

func (h *DummyHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("dummy", h.handle)
}

func (h *DummyHandler) handle(c *gin.Context) {
	response := gin.H{
		"success": true,
	}

	c.JSON(http.StatusOK, response)
}
