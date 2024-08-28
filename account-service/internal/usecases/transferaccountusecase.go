package usecases

import (
	"errors"
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/repositories"
)

type TransferAccountUseCaseInterface interface {
	Handle(fromNumber string, toNumber string, value int64) error
}

type TransferAccountUseCase struct {
	accountRepository repositories.AccountRepositoryInterface
}

func NewTransferAccountUseCase(accountRepository repositories.AccountRepositoryInterface) *TransferAccountUseCase {
	return &TransferAccountUseCase{
		accountRepository: accountRepository,
	}
}

func (us *TransferAccountUseCase) Handle(fromNumber string, toNumber string, value int64) error {
	fromAcc, err := us.accountRepository.GetAccountByNumber(fromNumber)
	if err != nil {
		slog.Error("Error getting from account by document", "error", err)
		return err
	}

	if fromAcc == nil {
		slog.Info("from account not found", "fromNumber", fromNumber)
		return errors.New("from account not found")
	}

	toAcc, err := us.accountRepository.GetAccountByNumber(toNumber)
	if err != nil {
		slog.Error("Error getting to account by document", "error", err)
		return err
	}

	if toAcc == nil {
		slog.Info("to account not found", "toNumber", toNumber)
		return errors.New("to account not found")
	}

	fromAcc.Transfer(value, toAcc)

	err = us.accountRepository.UpdateAccountBalance(toAcc)
	if err != nil {
		slog.Error("Error updating to account balance", "error", err)
		return err
	}

	err = us.accountRepository.UpdateAccountBalance(fromAcc)
	if err != nil {
		slog.Error("Error updating from account balance", "error", err)
	}

	return err
}
