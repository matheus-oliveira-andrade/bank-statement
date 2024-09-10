package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/infrastructure/broker"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/usecases"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/server/controllers"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/server/middleware"
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

	statementGenerationRepository := repositories.NewStatementGenerationRepository(repositories.NewDBConnection())
	accountRepository := repositories.NewAccountRepository(repositories.NewDBConnection())
	broker := broker.NewBroker(broker.BuildConnectionUrl())

	triggerStatementUseCase := usecases.NewTriggerStatementGenerationUseCase(statementGenerationRepository, accountRepository, broker)

	controllers.NewStatementController(triggerStatementUseCase).RegisterRoutes(v1Group)
}

func (s *APIServer) SetupMiddlewares() {
	s.Engine.Use(middleware.NewRequestLoggerMiddleware())
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
