package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/usecases"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/server/middleware"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/server/models"
)

type StatementController struct {
	triggerStatementGenerationUseCase usecases.TriggerStatementGenerationUseCaseInterface
}

func NewStatementController(triggerStatementGenerationUseCase usecases.TriggerStatementGenerationUseCaseInterface) *StatementController {
	return &StatementController{
		triggerStatementGenerationUseCase: triggerStatementGenerationUseCase,
	}
}

func (a *StatementController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/statement/:AccountNumber", middleware.NewAuthMiddleware("bankstatement"), a.triggerStatementGeneration)
}

func (c *StatementController) triggerStatementGeneration(ctx *gin.Context) {
	var req models.TriggerStatementGenerationRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		slog.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	triggerId, err := c.triggerStatementGenerationUseCase.Handle(req.AccountNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, models.NewTriggerStatementGenerationResponse(triggerId))
}
