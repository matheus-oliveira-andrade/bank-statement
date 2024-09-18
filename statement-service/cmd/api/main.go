package main

import (
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/configs"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/logger"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/server"
	"github.com/spf13/viper"
)

func main() {
	configs.InitConfigFile()

	logger.SetupLogger(viper.GetString("serviceName"))

	s := server.NewApiServer(viper.GetInt("port"))
	s.SetupMiddlewares()
	s.SetupRoutes()

	s.Start()
}
