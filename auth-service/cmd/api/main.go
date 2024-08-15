package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/matheus-oliveira-andrade/bank-statement/auth-service/internal/server"
	"github.com/spf13/viper"
)

type APIServer struct {
	port   int
	engine *gin.Engine
}

func NewApiServer(port int) *APIServer {
	return &APIServer{
		engine: gin.Default(),
		port:   port,
	}
}

func (s *APIServer) Setup() {
	baseGroup := s.engine.Group(viper.GetString("serviceBaseRoute"))
	server.NewHealthHandler().RegisterRoutes(baseGroup)

	v1Group := baseGroup.Group("v1")
	{
		server.NewDummyHandler().RegisterRoutes(v1Group)
		server.NewTokenHandler().RegisterRoutes(v1Group)
	}
}

func (s *APIServer) Start() {
	err := s.engine.Run(fmt.Sprint(":", s.port))
	if err != nil {
		log.Fatal(err)
	}
}

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

	s := NewApiServer(viper.GetInt("port"))
	s.Setup()

	s.Start()
}
