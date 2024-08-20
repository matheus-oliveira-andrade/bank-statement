package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func NewAuthMiddleware(requiredScope string) gin.HandlerFunc {
	return checkAuthHandle(requiredScope)
}

func checkAuthHandle(requiredScope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			slog.Info("not found token in header")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		secret := viper.GetString("authSettings.secret")
		if secret == "" {
			slog.Info("secret for JWT not found")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			slog.Error(err.Error())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !tokenParsed.Valid {
			slog.Info("token not valid")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !hasRequiredScope(tokenParsed.Claims.(jwt.MapClaims), requiredScope) {
			slog.Info("required scope to access not found")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

func hasRequiredScope(claims jwt.MapClaims, requiredScope string) bool {
	scopes := claims["scopes"].([]interface{})

	wasFound := false
	for _, scope := range scopes {
		if scope.(string) == requiredScope {
			wasFound = true
			break
		}
	}

	return wasFound
}
