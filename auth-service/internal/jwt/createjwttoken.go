package jwt

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type AuthManager struct {
}

func NewAuthManager() *AuthManager {
	return &AuthManager{}
}

func (*AuthManager) CreateJWTToken(accountNumber string) (string, error) {
	audience := viper.GetString("authSettings.audience")
	scopes := viper.GetStringSlice("authSettings.scopes")
	secret := viper.GetString("authSettings.secret")
	expirationHours := viper.GetInt("expirationHours")

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
		log.Println("Token not generated: ", err)
		return "", err
	}

	return tokenGenerated, nil
}
