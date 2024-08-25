package usecases

import (
	"errors"
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/repositories"
)

type DepositAccountUseCaseInterface interface {
	Handle(number string, value int64) error
}

type DepositAccountUseCase struct {
	accountRepository repositories.AccountRepositoryInterface
}

func NewDepositAccountUseCase(accountRepository repositories.AccountRepositoryInterface) *DepositAccountUseCase {
	return &DepositAccountUseCase{
		accountRepository: accountRepository,
	}
}

func (us *DepositAccountUseCase) Handle(number string, value int64) error {
	acc, err := us.accountRepository.GetAccountByNumber(number)
	if err != nil {
		slog.Error("Error getting account by document", "error", err)
		return err
	}

	if acc == nil {
		slog.Info("account not found", "number", number)
		return errors.New("account not found")
	}

	acc.Deposit(value)

	err = us.accountRepository.UpdateAccountBalance(acc)

	return err
}
