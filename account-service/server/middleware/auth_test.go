package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func generateTestJWTToken(accountNumber string, hoursToExpire int, scope string) string {
	secret := "123456"
	expirationTime := time.Now().Add(time.Hour * time.Duration(hoursToExpire)).Unix()

	claims := jwt.MapClaims{
		"exp": expirationTime,
		"sub": accountNumber,
		"aud": "test",
		"scopes": []string{
			scope,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtToken, _ := token.SignedString([]byte(secret))

	return jwtToken
}

func TestAuthHandler(t *testing.T) {
	testCases := []struct {
		name                   string
		requiredScope          string
		token                  string
		tokenSecret            string
		expectedHttpStatusCode int
	}{
		{
			name:                   "when dont have token should return unauthorized for empty token",
			requiredScope:          "",
			token:                  "",
			tokenSecret:            "",
			expectedHttpStatusCode: http.StatusForbidden,
		},
		{
			name:                   "when not set secret should return unauthorized",
			requiredScope:          "",
			token:                  generateTestJWTToken("123456", 1, ""),
			tokenSecret:            "",
			expectedHttpStatusCode: http.StatusInternalServerError,
		},
		{
			name:                   "when used an invalid signature should return unauthorized",
			requiredScope:          "",
			token:                  generateTestJWTToken("123456", -1, ""),
			tokenSecret:            "123321",
			expectedHttpStatusCode: http.StatusUnauthorized,
		},
		{
			name:                   "when have an expired token should return unauthorized",
			requiredScope:          "",
			token:                  generateTestJWTToken("123456", -1, ""),
			tokenSecret:            "123456",
			expectedHttpStatusCode: http.StatusUnauthorized,
		},
		{
			name:                   "when have a valid token and dont have scope should return unauthorized",
			requiredScope:          "ABC",
			token:                  generateTestJWTToken("123456", 1, "CBA"),
			tokenSecret:            "123456",
			expectedHttpStatusCode: http.StatusForbidden,
		},
		{
			name:                   "when have a valid token and have scope should return handle request",
			requiredScope:          "ABC",
			token:                  generateTestJWTToken("123456", 1, "ABC"),
			tokenSecret:            "123456",
			expectedHttpStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.ReleaseMode)

			defer viper.Reset()
			viper.Set("authSettings.secret", testCase.tokenSecret)

			router := gin.New()
			router.Use(NewCheckJWTTokenMiddleware(testCase.requiredScope))
			router.GET("/test", func(c *gin.Context) {
				c.String(http.StatusOK, "Hello, World!")
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			req.Header = map[string][]string{
				"Authorization": {testCase.token},
			}
			recorder := httptest.NewRecorder()

			// Act
			router.ServeHTTP(recorder, req)

			// Assert
			assert.Equal(t, testCase.expectedHttpStatusCode, recorder.Code)
		})
	}
}
