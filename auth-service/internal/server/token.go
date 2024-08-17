package server

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matheus-oliveira-andrade/bank-statement/auth-service/internal/jwt"
)

type AuthHandler struct {
}

func NewTokenHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("token", h.handle)
}

func (h *AuthHandler) handle(c *gin.Context) {
	tokenRequest := TokenRequest{}
	json.NewDecoder(c.Request.Body).Decode(&tokenRequest)

	token, err := jwt.NewAuthManager().CreateJWTToken(tokenRequest.AccountNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	response := gin.H{
		"token": token,
	}

	c.JSON(http.StatusOK, response)
}
