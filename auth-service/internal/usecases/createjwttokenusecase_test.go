package usecases

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestHandle(t *testing.T) {
	testsCase := []struct {
		name          string
		accountNumber string
	}{
		{
			name:          "check JWT creation",
			accountNumber: "123456-78",
		},
	}

	for _, tc := range testsCase {
		t.Run(tc.name, func(t *testing.T) {
			viper.Set("authSettings.secret", "secret")
			viper.Set("authSettings.audience", "webAPIs")
			viper.Set("authSettings.expirationHours", 4)
			viper.Set("authSettings.scopes", []string{
				"account",
				"bankstatement",
			})

			token, err := NewCreateJWTTokenUseCase().Handle(tc.accountNumber)

			assert.Nil(t, err)
			assert.NotNil(t, token)
		})
	}
}
