package usecases_mocks

import (
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/infrastructure/broker"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/shared/events"
	"github.com/stretchr/testify/mock"
)

type MockBroker struct {
	mock.Mock
}

func (m *MockBroker) Produce(eventPublish *events.EventPublish, configs *broker.ProduceConfigs) error {
	args := m.Called(eventPublish, configs)
	return args.Error(0)
}
