package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/matheus-oliveira-andrade/bank-statement/account-service/server"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func initTestConfigFile() {
	viper.AddConfigPath("../../configs")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	environment := viper.GetString("environment")

	viper.SetConfigName(fmt.Sprint("configs", ".", environment))
	viper.SetConfigType("json")
	viper.AddConfigPath("configs")

	err = viper.MergeInConfig()
	if err != nil {
		panic(err)
	}
}
func TestApiServerStart(t *testing.T) {
	gin.SetMode(gin.TestMode)

	initTestConfigFile()

	server := server.NewApiServer(8080)
	server.SetupMiddlewares()
	server.SetupRoutes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/account/health", nil)
	server.Engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
