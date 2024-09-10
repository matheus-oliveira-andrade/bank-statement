package handlersmock

import (
	"github.com/stretchr/testify/mock"
)

type MockGenerateDocumentApi struct {
	mock.Mock
}

func (m *MockGenerateDocumentApi) GenerateFromHtml(html string) (string, error) {
	args := m.Called(html)
	return args.String(0), args.Error(1)
}
