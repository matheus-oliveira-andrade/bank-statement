package eventhandlers

import (
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/shared/events"
)

type TransferReceivedHandlerInterface interface {
	Handler(event events.TransferReceived)
}

type TransferReceivedHandler struct {
	accountRepository  repositories.AccountRepositoryInterface
	movementRepository repositories.MovementRepositoryInterface
}

func NewTransferReceivedHandler(accountRepository repositories.AccountRepositoryInterface, movementRepository repositories.MovementRepositoryInterface) TransferReceivedHandlerInterface {
	return &TransferReceivedHandler{
		accountRepository:  accountRepository,
		movementRepository: movementRepository,
	}
}

func (h *TransferReceivedHandler) Handler(event events.TransferReceived) {
	slog.Info("handling transfer received", "number", event.FromNumber)

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

	movement := domain.NewTransferReceivedMovement(event.FromNumber, event.ToNumber, event.Value)
	err = h.movementRepository.CreateMovement(movement)
	if err != nil {
		slog.Error("error creating movement", "error", err, "number", event.FromNumber)
		return
	}

	slog.Info("transfer received account updated", "number", event.FromNumber)
}
