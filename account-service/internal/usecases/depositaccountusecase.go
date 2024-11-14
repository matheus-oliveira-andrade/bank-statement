package usecases

import (
	"errors"
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/infrastructure/broker"
	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/bank-statement/account-service/shared/events"
)

type DepositAccountUseCaseInterface interface {
	Handle(number string, value int64, idempotencyKey string) error
}

type DepositAccountUseCase struct {
	accountRepository         repositories.AccountRepositoryInterface
	broker                    broker.BrokerInterface
	idempotencyKeysRepository repositories.IdempotencyKeysRepositoryInterface
}

func NewDepositAccountUseCase(
	accountRepository repositories.AccountRepositoryInterface,
	broker broker.BrokerInterface,
	idempotencyKeysRepository repositories.IdempotencyKeysRepositoryInterface) *DepositAccountUseCase {
	return &DepositAccountUseCase{
		accountRepository:         accountRepository,
		broker:                    broker,
		idempotencyKeysRepository: idempotencyKeysRepository,
	}
}

func (us *DepositAccountUseCase) Handle(number string, value int64, idempotencyKey string) error {
	hasKey, err := us.idempotencyKeysRepository.HasKey(idempotencyKey)
	if err != nil {
		slog.Error("Error getting idempotencyKey", "error", err)
		return err
	}

	if hasKey {
		slog.Error("idempotency key already processed", "idempotencyKey", idempotencyKey)
		return errors.New("idempotency key already processed")
	}

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
	if err != nil {
		slog.Error("error updating account balance", "error", err)
		return err
	}

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

	err = us.idempotencyKeysRepository.CreateKey(idempotencyKey)
	if err != nil {
		slog.Error("error saving idempotency key used", "error", err, "idempotencyKey", idempotencyKey)
		return err
	}

	slog.Info("Deposit created", "accountNumber", acc.Number, "value", value, "idempotencyKey", idempotencyKey)

	return err
}
