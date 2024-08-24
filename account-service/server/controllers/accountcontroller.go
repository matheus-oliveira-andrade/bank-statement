package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/usecases"
	"github.com/matheus-oliveira-andrade/bank-statement/account-service/server/middleware"
	"github.com/matheus-oliveira-andrade/bank-statement/account-service/server/models"
)

type AccountController struct {
	createAccountUseCase usecases.CreateAccountUseCaseInterface
}

func NewAccountController(createAccountUseCase usecases.CreateAccountUseCaseInterface) *AccountController {
	return &AccountController{
		createAccountUseCase: createAccountUseCase,
	}
}

func (a *AccountController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/account", middleware.NewAuthMiddleware("account"), a.handle)
}

func (c *AccountController) handle(ctx *gin.Context) {
	req := models.CreateAccountRequest{}
	err := json.NewDecoder(ctx.Request.Body).Decode(&req)
	if err != nil {
		slog.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	number, err := c.createAccountUseCase.Handle(req.Document, req.Name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, models.NewCreateAccountResponse(number))
}
