package main

import (
	"fmt"

	"github.com/matheus-oliveira-andrade/bank-statement/auth-service/internal/logger"
	"github.com/matheus-oliveira-andrade/bank-statement/auth-service/server"
	"github.com/spf13/viper"
)

func initConfigFile() {
	viper.AddConfigPath("configs")
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

func main() {
	initConfigFile()

	logger.SetupLogger(viper.GetString("serviceName"))

	s := server.NewApiServer(viper.GetInt("port"))
	s.SetupMiddlewares()
	s.SetupRoutes()

	s.Start()
}
