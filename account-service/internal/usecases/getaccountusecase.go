package usecases

import (
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/domain"
	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/repositories"
)

type GetAccountUseCaseInterface interface {
	Handle(number string) (*domain.Account, error)
}

type GetAccountUseCase struct {
	accountRepository repositories.AccountRepositoryInterface
}

func NewGetAccountUseCase(accountRepository repositories.AccountRepositoryInterface) *GetAccountUseCase {
	return &GetAccountUseCase{
		accountRepository: accountRepository,
	}
}

func (us *GetAccountUseCase) Handle(number string) (*domain.Account, error) {
	acc, err := us.accountRepository.GetAccountByNumber(number)
	if err != nil {
		slog.Info("error get account by number", "error", err)
		return nil, err
	}

	slog.Info("searched account by number", "number", number, "acc", acc)
	return acc, nil
}
