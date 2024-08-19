package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/matheus-oliveira-andrade/bank-statement/auth-service/internal/usecases"
	"github.com/matheus-oliveira-andrade/bank-statement/auth-service/server/controllers"
	"github.com/matheus-oliveira-andrade/bank-statement/auth-service/server/middleware"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

type APIServer struct {
	port   int
	Engine *gin.Engine
}

func NewApiServer(port int) *APIServer {
	return &APIServer{
		Engine: gin.New(),
		port:   port,
	}
}

func (s *APIServer) SetupRoutes() {
	baseGroup := s.Engine.Group(viper.GetString("serviceBaseRoute"))
	controllers.NewHealthController().RegisterRoutes(baseGroup)

	v1Group := baseGroup.Group("v1")
	createTokenUseCase := usecases.NewCreateJWTTokenUseCase()
	controllers.NewAuthController(createTokenUseCase).RegisterRoutes(v1Group)
}

func (s *APIServer) SetupMiddlewares() {
	s.Engine.Use(middleware.DefaultStructuredLogger())
	s.Engine.Use(gin.Recovery())
}

func (s *APIServer) Start() {
	err := s.Engine.Run(fmt.Sprint(":", s.port))
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	slog.Info("server started", "port", s.port)
}
