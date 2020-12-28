package internal

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
	"webhook-server-test/internal/service"
)

type Server struct {
	*http.Server
	l *zap.Logger
}

func NewServer(addr string, logger *zap.Logger, svc *service.Service) *Server {
	mux := http.NewServeMux()
	mux.Handle("/webhooks", http.HandlerFunc(svc.CreateWebhooksHandler))
	mux.Handle("/fire-webhooks", http.HandlerFunc(svc.FireWebhooksHandler))

	return &Server{
		l: logger,
		Server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (s *Server) StartWithGracefulShutdown() {
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		s.l.Info("server is shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		s.Server.SetKeepAlivesEnabled(false)
		if err := s.Server.Shutdown(ctx); err != nil {
			s.l.Error(fmt.Sprintf("graceful shutdown failed: %v\n", err))
		}
		close(done)
	}()

	if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.l.Fatal(fmt.Sprintf("server start failed: %v\n", err))
	}

	<-done
	s.l.Info("server stopped")
}
