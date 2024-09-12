package usecases

import (
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type CreateJWTTokenUseCaseInterface interface {
	Handle() (string, error)
}

type CreateJWTTokenUseCase struct {
}

func NewCreateJWTTokenUseCase() *CreateJWTTokenUseCase {
	return &CreateJWTTokenUseCase{}
}

func (*CreateJWTTokenUseCase) Handle() (string, error) {
	slog.Info("Creating JWT token")

	audience := viper.GetString("authSettings.audience")
	scopes := viper.GetStringSlice("authSettings.scopes")
	secret := viper.GetString("authSettings.secret")
	expirationHours := viper.GetInt("authSettings.expirationHours")

	expirationTime := time.Now().Add(time.Hour * time.Duration(expirationHours)).Unix()

	claims := jwt.MapClaims{
		"exp":    expirationTime,
		"sub":    uuid.New().String(),
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

	slog.Info("Token created")

	return tokenGenerated, nil
}
