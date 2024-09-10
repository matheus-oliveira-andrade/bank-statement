package eventhandlers

import (
	"fmt"
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	documentgenerator "github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/infrastructure/documentgenerator"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/infrastructure/templatecompiler"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/shared/events"
)

type StatementGenerationRequestedHandlerInterface interface {
	Handle(event events.StatementGenerationRequested)
}

type StatementGenerationRequestedHandler struct {
	accountRepository             repositories.AccountRepositoryInterface
	statementGenerationRepository repositories.StatementGenerationRepositoryInterface
	movementRepository            repositories.MovementRepositoryInterface
	documentGeneratorApi          documentgenerator.GenerateDocumentApiInterface
	templateCompiler              templatecompiler.TemplateCompileInterface
}

func NewStatementGenerationRequestedHandler(
	accountRepository repositories.AccountRepositoryInterface,
	statementGenerationRepository repositories.StatementGenerationRepositoryInterface,
	movementRepository repositories.MovementRepositoryInterface,
	documentGeneratorApi documentgenerator.GenerateDocumentApiInterface,
	templateCompiler templatecompiler.TemplateCompileInterface,
) StatementGenerationRequestedHandlerInterface {
	return &StatementGenerationRequestedHandler{
		accountRepository:             accountRepository,
		statementGenerationRepository: statementGenerationRepository,
		movementRepository:            movementRepository,
		documentGeneratorApi:          documentGeneratorApi,
		templateCompiler:              templateCompiler,
	}
}

func (us *StatementGenerationRequestedHandler) Handle(event events.StatementGenerationRequested) {
	slog.Info("handling account created", "number", event.AccountNumber)

	statementGeneration, err := us.statementGenerationRepository.GetStatementGeneration(event.AccountNumber)
	if err != nil {
		slog.Error("error generating statement", "error", err)
		return
	}

	acc, err := us.accountRepository.GetAccountByNumber(event.AccountNumber)
	if err != nil {
		slog.Error("error getting account", "error", err)
		us.UpdateStatementGenerationError(statementGeneration, err)
		return
	}

	if acc == nil {
		slog.Error("account not found", "number", event.AccountNumber)
		return
	}

	movements, err := us.movementRepository.GetMovements(event.AccountNumber)
	if err != nil {
		slog.Error("error getting movements", "error", err)
		us.UpdateStatementGenerationError(statementGeneration, err)
		return
	}

	if movements == nil {
		slog.Error("movements not found", "number", event.AccountNumber)
		return
	}

	parameters := us.NewStatementGenerationReportParameter(acc, movements, statementGeneration)

	templateCompiled, err := us.templateCompiler.Compile(parameters)
	if err != nil {
		us.UpdateStatementGenerationError(statementGeneration, err)
		return
	}

	report, err := us.documentGeneratorApi.GenerateFromHtml(templateCompiled)
	if err != nil {
		slog.Error("error generating document", "error", err)
		us.UpdateStatementGenerationError(statementGeneration, err)
		return
	}

	statementGeneration.SetAsGenerated(report)

	err = us.statementGenerationRepository.UpdateStatementGeneration(statementGeneration)
	if err != nil {
		slog.Error("error updating statement generation", "error", err)
		return
	}

}

func (us *StatementGenerationRequestedHandler) UpdateStatementGenerationError(sg *domain.StatementGeneration, err error) {
	sg.SetAsGeneratedWithError(err)

	err = us.statementGenerationRepository.UpdateStatementGeneration(sg)
	if err != nil {
		slog.Error("error updating statement generation", "error", err)
		return
	}
}

func (us *StatementGenerationRequestedHandler) NewStatementGenerationReportParameter(
	acc *domain.Account,
	movements *[]domain.Movement,
	sg *domain.StatementGeneration) *domain.StatementGenerationReportParameter {

	reportParameter := domain.StatementGenerationReportParameter{
		Document:     acc.Document,
		CustomerName: acc.Name,
		Movements:    []domain.MovementReportParameter{},
	}

	for _, movement := range *movements {
		movementType := ""
		if movement.Type == string(domain.In) {
			movementType = "Entrada"
		} else {
			movementType = "Saída"
		}

		destinationAccount := ""
		if movement.ToAccountNumber == "" {
			destinationAccount = " - "
		} else {
			destinationAccount = movement.ToAccountNumber
		}

		movementParameter := domain.MovementReportParameter{
			CreatedAt:          movement.CreatedAt.Format("2006-01-02 15:04:05"),
			Type:               movementType,
			DestinationAccount: destinationAccount,
			Amount:             fmt.Sprintf("R$ %.2f", float32(movement.Value/100)),
		}

		reportParameter.Movements = append(reportParameter.Movements, movementParameter)
	}

	return &reportParameter
}
