package httpService

import (
	"context"
	"fmt"
	"github.com/obsidian-project-plugin/auth-service/internal/config"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	mux    *http.ServeMux
	config config.ServerConfig
}

func NewServer(cfg config.ServerConfig) *Server {
	mux := http.NewServeMux()

	return &Server{
		server: &http.Server{},
		mux:    mux,
		config: cfg,
	}
}

func (s *Server) Mux() *http.ServeMux {
	return s.mux
}

func (s *Server) Address() string {
	return fmt.Sprintf(":%s", s.config.HTTPPort)
}

func (s *Server) ListenAndServe() error {
	s.server.Addr = s.Address()
	s.server.Handler = s.Mux()
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}
