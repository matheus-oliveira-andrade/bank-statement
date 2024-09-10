package handlersmock

import (
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockTemplateCompiler struct {
	mock.Mock
}

func (m *MockTemplateCompiler) Compile(parameters *domain.StatementGenerationReportParameter) (string, error) {
	args := m.Called(parameters)
	return args.String(0), args.Error(1)
}
