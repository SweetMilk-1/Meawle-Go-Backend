package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"meawle/internal/config"
)

// Server представляет HTTP сервер приложения
type Server struct {
	*http.Server
	Logger *log.Logger
}

// NewServer создает новый HTTP сервер
func NewServer(cfg *config.Config, handler http.Handler, logger *log.Logger) *Server {
	return &Server{
		Server: &http.Server{
			Addr:         cfg.Port,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Logger: logger,
	}
}

// Start запускает сервер в горутине
func (s *Server) Start() {
	go func() {
		s.Logger.Printf("Server starting on port %s", s.Addr)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Fatal("Server failed to start:", err)
		}
	}()
}

// Shutdown выполняет graceful shutdown сервера
func (s *Server) Shutdown() {
	s.Logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Server.Shutdown(ctx); err != nil {
		s.Logger.Fatal("Server forced to shutdown:", err)
	}

	s.Logger.Println("Server exited")
}

// WaitForShutdown ожидает сигналов для graceful shutdown
func WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}