package eventhandlers

import (
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/shared/events"
)

type FundsDepositedHandlerInterface interface {
	Handler(event events.FundsDeposited)
}

type FundsDepositedHandler struct {
	accountRepository  repositories.AccountRepositoryInterface
	movementRepository repositories.MovementRepositoryInterface
}

func NewFundsDepositedHandler(accountRepository repositories.AccountRepositoryInterface, movementRepository repositories.MovementRepositoryInterface) FundsDepositedHandlerInterface {
	return &FundsDepositedHandler{
		accountRepository:  accountRepository,
		movementRepository: movementRepository,
	}
}

func (h *FundsDepositedHandler) Handler(event events.FundsDeposited) {
	slog.Info("handling funds deposited", "number", event.Number)

	acc, err := h.accountRepository.GetAccountByNumber(event.Number)
	if err != nil {
		slog.Error("error getting account", "error", err)
		return
	}

	if acc == nil {
		slog.Error("account not found", "number", event.Number)
		return
	}

	acc.Balance += event.Value

	err = h.accountRepository.UpdateAccountBalance(acc)
	if err != nil {
		slog.Error("error updating account balance", "error", err, "number", event.Number)
		return
	}

	movement := domain.NewDepositedFundsMovement(event.Number, event.Value)
	err = h.movementRepository.CreateMovement(movement)
	if err != nil {
		slog.Error("error creating movement", "error", err, "number", event.Number)
		return
	}

	slog.Info("funds deposited account updated", "number", event.Number)
}
