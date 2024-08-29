package main

import (
	"fmt"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/logger"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/server"
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

	logger.SetupLogger()

	s := server.NewApiServer(viper.GetInt("port"))
	s.SetupMiddlewares()
	s.SetupRoutes()

	s.Start()
}
