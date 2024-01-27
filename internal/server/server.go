package server

import (
	"log/slog"

	"github.com/Cataloft/user-service/internal/config"
	"github.com/Cataloft/user-service/internal/handlers/user"
	"github.com/Cataloft/user-service/internal/middlewares"
	"github.com/Cataloft/user-service/internal/storage"
	"github.com/gin-gonic/gin"
	requestid "github.com/sumit-tembe/gin-requestid"
)

type Server struct {
	router  *gin.Engine
	storage *storage.Storage
	config  *config.Config
	logger  *slog.Logger
}

func New(db *storage.Storage, cfg *config.Config, logger *slog.Logger) *Server {
	router := gin.New()

	router.Use(middlewares.LogMiddleware(logger))
	router.Use(gin.Recovery())
	router.Use(requestid.RequestID(nil))

	return &Server{
		router:  router,
		storage: db,
		config:  cfg,
		logger:  logger,
	}
}

func (s *Server) initHandlers() {
	s.router.Use(middlewares.PaginationMiddleware())
	s.router.POST("/users", user.EnrichFIO(s.storage, s.logger))
	s.router.GET("/users", user.GetList(s.storage, s.logger))

	s.router.DELETE("/delete/:id", user.Delete(s.storage, s.logger))
	s.router.PATCH("/update/:id", user.Update(s.storage, s.logger))
}

func (s *Server) Start() error {
	s.initHandlers()

	return s.router.Run(s.config.ServerAddr)
}
