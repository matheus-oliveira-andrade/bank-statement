package eventhandlers

import (
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/shared/events"
)

type TransferRealizedHandlerInterface interface {
	Handler(event events.TransferRealized)
}

type TransferRealizedHandler struct {
	accountRepository  repositories.AccountRepositoryInterface
	movementRepository repositories.MovementRepositoryInterface
}

func NewTransferRealizedHandler(accountRepository repositories.AccountRepositoryInterface, movementRepository repositories.MovementRepositoryInterface) TransferRealizedHandlerInterface {
	return &TransferRealizedHandler{
		accountRepository:  accountRepository,
		movementRepository: movementRepository,
	}
}

func (h *TransferRealizedHandler) Handler(event events.TransferRealized) {
	slog.Info("handling transfer realized", "number", event.FromNumber)

	acc, err := h.accountRepository.GetAccountByNumber(event.FromNumber)
	if err != nil {
		slog.Error("error getting account", "error", err)
		return
	}

	if acc == nil {
		slog.Error("account not found", "number", event.FromNumber)
		return
	}

	acc.Balance = event.Balance

	err = h.accountRepository.UpdateAccountBalance(acc)
	if err != nil {
		slog.Error("error updating account balance", "error", err, "number", event.FromNumber)
		return
	}

	movement := domain.NewTransferRealizedMovement(event.FromNumber, event.ToNumber, event.Value)
	err = h.movementRepository.CreateMovement(movement)
	if err != nil {
		slog.Error("error creating movement", "error", err, "number", event.FromNumber)
		return
	}

	slog.Info("transfer realized account updated", "number", event.FromNumber)
}
