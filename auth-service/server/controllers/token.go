package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matheus-oliveira-andrade/bank-statement/auth-service/internal/usecases"
	"github.com/matheus-oliveira-andrade/bank-statement/auth-service/server/models"
)

type AuthController struct {
	CreateTokenUseCase usecases.CreateJWTTokenUseCaseInterface
}

func NewAuthController(createTokenUseCase usecases.CreateJWTTokenUseCaseInterface) *AuthController {
	return &AuthController{
		CreateTokenUseCase: createTokenUseCase,
	}
}

func (controller *AuthController) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("token", controller.CreateToken)
}

func (controller *AuthController) CreateToken(ctx *gin.Context) {
	tokenRequest := models.TokenRequest{}
	json.NewDecoder(ctx.Request.Body).Decode(&tokenRequest)

	token, err := controller.CreateTokenUseCase.Handle(tokenRequest.AccountNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	response := gin.H{
		"token": token,
	}

	ctx.JSON(http.StatusOK, response)
}
