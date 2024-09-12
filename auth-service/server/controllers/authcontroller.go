package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matheus-oliveira-andrade/bank-statement/auth-service/internal/usecases"
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
	token, err := controller.CreateTokenUseCase.Handle()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	response := gin.H{
		"token": token,
	}

	ctx.JSON(http.StatusOK, response)
}
