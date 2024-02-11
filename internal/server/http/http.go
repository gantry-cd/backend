package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gantrycd/backend/internal/server"
)

type httpServer struct {
	port            string
	host            string
	shutdownTimeout time.Duration

	l *slog.Logger

	srv *http.Server
}

// Option はサーバーのオプションを設定するための関数です。
type Option func(*httpServer)

// WithPort はポート番号を設定するオプションです。
func WithPort(port string) Option {
	return func(s *httpServer) {
		s.port = port
	}
}

// WithHost はホスト名を設定するオプションです。
func WithHost(host string) Option {
	return func(s *httpServer) {
		s.host = host
	}
}

// WithLogger はロガーを設定するオプションです。
func WithLogger(l *slog.Logger) Option {
	return func(s *httpServer) {
		s.l = l
	}
}

// WithShutdownTimeout はシャットダウンタイムアウトを設定するオプションです。
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(s *httpServer) {
		s.shutdownTimeout = timeout
	}
}

// New はサーバーを生成します。
func New(handler http.Handler, opts ...Option) server.Server {
	s := &httpServer{
		port:            "8080",
		host:            "localhost",
		shutdownTimeout: server.DefaultShutdownTimeout,
		l:               slog.New(slog.NewTextHandler(os.Stderr, nil)).WithGroup("server"),
	}

	for _, opt := range opts {
		opt(s)
	}

	s.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.host, s.port),
		Handler: handler,
	}

	return s
}

// Run はサーバーを起動します。
func (s *httpServer) Run() error {
	s.l.Info(fmt.Sprintf("server starting at %s", s.srv.Addr))
	return s.srv.ListenAndServe()
}

// Shutdown はサーバーを停止します。
func (s *httpServer) Shutdown(ctx context.Context) error {
	s.l.Info("server shutdown ...")
	return s.srv.Shutdown(ctx)
}

// RunWithGracefulShutdown はgraceful shutdownを行うサーバーを起動します。
func (s *httpServer) RunWithGracefulShutdown() {
	go func() {
		if err := s.Run(); err != nil && err != http.ErrServerClosed {
			s.l.Error(fmt.Sprintf("Listen And Serve error : %s", err.Error()))
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	if err := s.Shutdown(ctx); err != nil {
		s.l.Error(fmt.Sprintf("server shutdown error : %s", err.Error()))
	}

	s.l.Info("server shutdown completed")
}
