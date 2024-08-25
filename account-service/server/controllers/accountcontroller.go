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
	getAccountUseCase    usecases.GetAccountUseCaseInterface
}

func NewAccountController(createAccountUseCase usecases.CreateAccountUseCaseInterface, getAccountUseCase usecases.GetAccountUseCaseInterface) *AccountController {
	return &AccountController{
		createAccountUseCase: createAccountUseCase,
		getAccountUseCase:    getAccountUseCase,
	}
}

func (a *AccountController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/account", middleware.NewAuthMiddleware("account"), a.createAccountHandler)
	router.GET("/account/:number", middleware.NewAuthMiddleware("account"), a.getAccountHandler)
}

func (c *AccountController) createAccountHandler(ctx *gin.Context) {
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

func (c *AccountController) getAccountHandler(ctx *gin.Context) {
	var req models.GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		slog.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	acc, err := c.getAccountUseCase.Handle(req.Number)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})

		return
	}

	if acc == nil {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	ctx.JSON(http.StatusOK, models.NewGetAccountResponse(acc))
}
