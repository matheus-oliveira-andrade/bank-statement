package usecases

import (
	"errors"
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/infrastructure/broker"
	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/bank-statement/account-service/shared/events"
)

type DepositAccountUseCaseInterface interface {
	Handle(number string, value int64) error
}

type DepositAccountUseCase struct {
	accountRepository repositories.AccountRepositoryInterface
	broker            broker.BrokerInterface
}

func NewDepositAccountUseCase(accountRepository repositories.AccountRepositoryInterface, broker broker.BrokerInterface) *DepositAccountUseCase {
	return &DepositAccountUseCase{
		accountRepository: accountRepository,
		broker:            broker,
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

	event, err := events.NewEventPublish(events.NewFundsDeposited(acc.Number, value))
	if err != nil {
		slog.Error("error creating funds deposited event", "error", err)
		return err
	}

	err = us.broker.Produce(event, &broker.ProduceConfigs{Topic: "account"})
	if err != nil {
		slog.Error("error producing event", "error", err)
		return err
	}

	return err
}
