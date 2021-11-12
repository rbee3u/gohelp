package confhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	kernel *gin.Engine `env:"-"`
	Port   int32       `env:""`
}

type ServerOption func(*Server)

func WithPort(port int32) ServerOption {
	return func(s *Server) {
		s.Port = port
	}
}

func New(opts ...ServerOption) (*Server, error) {
	s := &Server{}

	if err := s.SetDefaults(); err != nil {
		return nil, fmt.Errorf("failed to set defaults: %w", err)
	}

	for _, opt := range opts {
		opt(s)
	}

	if err := s.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize: %w", err)
	}

	return s, nil
}

func (s *Server) SetDefaults() error {
	if s.Port == 0 {
		s.Port = 80
	}

	return nil
}

func (s *Server) Initialize() error {
	gin.SetMode(gin.ReleaseMode)

	s.kernel = gin.New()

	return nil
}

func (s *Server) Kernel() *gin.Engine {
	return s.kernel
}

func (s *Server) Serve(ctx context.Context) error {
	addr := fmt.Sprintf(":%v", s.Port)
	srv := &http.Server{Handler: s.kernel, Addr: addr}

	go func() { <-ctx.Done(); _ = srv.Shutdown(ctx) }()

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to listen and serve: %w", err)
	}

	return nil
}
