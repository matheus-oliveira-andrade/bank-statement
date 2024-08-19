package usecases

import (
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type CreateJWTTokenUseCaseInterface interface {
	Handle(accountNumber string) (string, error)
}

type CreateJWTTokenUseCase struct {
}

func NewCreateJWTTokenUseCase() *CreateJWTTokenUseCase {
	return &CreateJWTTokenUseCase{}
}

func (*CreateJWTTokenUseCase) Handle(accountNumber string) (string, error) {
	slog.Info("Creating JWT token", "accountNumber", accountNumber)

	audience := viper.GetString("authSettings.audience")
	scopes := viper.GetStringSlice("authSettings.scopes")
	secret := viper.GetString("authSettings.secret")
	expirationHours := viper.GetInt("authSettings.expirationHours")

	expirationTime := time.Now().Add(time.Hour * time.Duration(expirationHours)).Unix()

	claims := jwt.MapClaims{
		"exp":    expirationTime,
		"sub":    accountNumber,
		"aud":    audience,
		"scopes": scopes,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims)

	tokenGenerated, err := token.SignedString([]byte(secret))
	if err != nil {
		slog.Error("Token not generated: ", "err", err.Error())
		return "", err
	}

	slog.Info("Token created", "accountNumber", accountNumber)

	return tokenGenerated, nil
}