package main

import (
	"fmt"
	"github.com/matheus-oliveira-andrade/bank-statement/auth-service/internal/server"
	"log"

	"github.com/gin-gonic/gin"
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

func (s *APIServer) Start() {
	baseGroup := s.engine.Group("auth")
	server.NewHealthHandler().RegisterRoutes(baseGroup)

	v1Group := baseGroup.Group("v1")
	{
		server.NewDummyHandler().RegisterRoutes(v1Group)
	}

	err := s.engine.Run(fmt.Sprint(":", s.port))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	NewApiServer(8080).Start()
}
