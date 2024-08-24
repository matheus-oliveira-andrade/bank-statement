package usecases

import (
	"errors"
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/domain"
	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/repositories"
)

type CreateAccountUseCaseInterface interface {
	Handle(document string, name string) (string, error)
}

type CreateAccountUseCase struct {
	accountRepository repositories.AccountRepositoryInterface
}

func NewCreateAccountUseCase(accountRepository repositories.AccountRepositoryInterface) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		accountRepository: accountRepository,
	}
}

func (us *CreateAccountUseCase) Handle(document string, name string) (string, error) {
	acc, err := us.accountRepository.GetAccountByDocument(document)
	if err != nil {
		slog.Error("Error getting account by document", "error", err)
		return "", err
	}

	if acc != nil {
		slog.Info("document in use by another account", "document", document)
		return "", errors.New("document in use by another account")
	}

	number, err := us.accountRepository.GetNextAccountNumber()
	if err != nil {
		slog.Error("error get next account number", "error", err)
		return "", err
	}

	account := domain.NewAccount(number, document, name)
	err = account.Validate()
	if err != nil {
		slog.Error("Invalid account", "error", err)
		return "", err
	}

	id, err := us.accountRepository.CreateAccount(account)
	if err != nil {
		slog.Error("error creating account", "error", err)
		return "", err
	}

	slog.Info("account created", "id", id, "document", document, "name", name)
	return id, nil
}
