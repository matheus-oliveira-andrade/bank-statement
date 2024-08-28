package usecases

import (
	"errors"
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/infrastructure/broker"
	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/bank-statement/account-service/shared/events"
)

type TransferAccountUseCaseInterface interface {
	Handle(fromNumber string, toNumber string, value int64) error
}

type TransferAccountUseCase struct {
	accountRepository repositories.AccountRepositoryInterface
	broker            broker.BrokerInterface
}

func NewTransferAccountUseCase(accountRepository repositories.AccountRepositoryInterface, broker broker.BrokerInterface) *TransferAccountUseCase {
	return &TransferAccountUseCase{
		accountRepository: accountRepository,
		broker:            broker,
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
		return err
	}

	err = us.produceEventTransferRealized(fromAcc.Number, toAcc.Number, value, fromAcc.Balance)
	if err != nil {
		slog.Error("error producing event transfer realized", "error", err)
		return err
	}

	err = us.produceEventTransferReceived(toAcc.Number, fromAcc.Number, value, fromAcc.Balance)
	if err != nil {
		slog.Error("error producing event transfer received", "error", err)
		return err
	}

	return err
}

func (us *TransferAccountUseCase) produceEventTransferRealized(fromNumber string, toNumber string, value int64, fromBalance int64) error {
	transferRealizedEvent, err := events.NewEventPublish(events.NewTransferRealized(fromNumber, toNumber, value, fromBalance))
	if err != nil {
		return err
	}

	err = us.broker.Produce(transferRealizedEvent, &broker.ProduceConfigs{Topic: "account"})
	if err != nil {
		return err
	}

	return nil
}

func (us *TransferAccountUseCase) produceEventTransferReceived(toNumber string, fromNumber string, value int64, toBalance int64) error {
	transferReceivedEvent, err := events.NewEventPublish(events.NewTransferReceived(toNumber, fromNumber, value, toBalance))
	if err != nil {
		return err
	}

	err = us.broker.Produce(transferReceivedEvent, &broker.ProduceConfigs{Topic: "account"})
	if err != nil {
		return err
	}

	return nil
}
