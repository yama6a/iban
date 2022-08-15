package http

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// Controller an interface that defines a struct that can add routes to the http server via the default http package
type Controller interface {
	SetupRoutes()
}

// Server a web server that can serve multiple controllers.
type Server struct {
	port        uint
	logger      *zap.Logger
	controllers []Controller
}

// NewHttpServer creates a new http server.
func NewHttpServer(port uint, logger *zap.Logger, controllers []Controller) *Server {
	server := &Server{
		port:        port,
		logger:      logger,
		controllers: controllers,
	}

	server.setupRoutes()

	return server
}

// setupRoutes adds all the controllers routes to the http server.
func (s *Server) setupRoutes() {
	for _, controller := range s.controllers {
		controller.SetupRoutes()
	}
}

// Run starts the http server.
func (s *Server) Run() error {
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
	return fmt.Errorf("http server terminated: %w", err)
}
