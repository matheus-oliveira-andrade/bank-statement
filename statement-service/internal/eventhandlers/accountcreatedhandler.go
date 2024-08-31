package eventhandlers

import (
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/shared/events"
)

type AccountCreatedHandlerInterface interface {
	Handler(event events.AccountCreated)
}

type AccountCreatedHandler struct {
	accountRepository repositories.AccountRepositoryInterface
}

func NewAccountCreatedHandler(accountRepository repositories.AccountRepositoryInterface) AccountCreatedHandlerInterface {
	return &AccountCreatedHandler{
		accountRepository: accountRepository,
	}
}

func (h *AccountCreatedHandler) Handler(event events.AccountCreated) {
	slog.Info("handling account created", "number", event.Number)

	acc := domain.NewAccount(event.Number, event.Document, event.Name)

	err := h.accountRepository.CreateAccount(acc)

	if err != nil {
		slog.Error("error creating account", "error", err)
		return
	}

	slog.Info("account created", "number", event.Number)
}
