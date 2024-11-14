package usecases_mock

import (
	"github.com/stretchr/testify/mock"
)

type MockIdempotencyRepository struct {
	mock.Mock
}

func (m *MockIdempotencyRepository) HasKey(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *MockIdempotencyRepository) CreateKey(key string) error {
	args := m.Called(key)
	return args.Error(0)
}
