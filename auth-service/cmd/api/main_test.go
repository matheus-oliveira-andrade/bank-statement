package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/matheus-oliveira-andrade/bank-statement/auth-service/server"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestApiServerStart(t *testing.T) {
	gin.SetMode(gin.TestMode)

	viper.Set("port", 8080)
	viper.Set("serviceBaseRoute", "auth")

	server := server.NewApiServer(8080)
	server.SetupMiddlewares()
	server.SetupRoutes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/auth/health", nil)
	server.Engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
