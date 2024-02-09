package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"log/slog"
)

// DefaultShutdownTimeout はデフォルトのシャットダウンタイムアウトです。
const DefaultShutdownTimeout time.Duration = 10

type server struct {
	port            string
	host            string
	shutdownTimeout time.Duration

	l *slog.Logger

	srv *http.Server
}

// Option はサーバーのオプションを設定するための関数です。
type Option func(*server)

// WithPort はポート番号を設定するオプションです。
func WithPort(port string) Option {
	return func(s *server) {
		s.port = port
	}
}

// WithHost はホスト名を設定するオプションです。
func WithHost(host string) Option {
	return func(s *server) {
		s.host = host
	}
}

// WithLogger はロガーを設定するオプションです。
func WithLogger(l *slog.Logger) Option {
	return func(s *server) {
		s.l = l
	}
}

// WithShutdownTimeout はシャットダウンタイムアウトを設定するオプションです。
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(s *server) {
		s.shutdownTimeout = timeout
	}
}

// Server はサーバーを表すインターフェースです。
type Server interface {
	Run() error
	Shutdown(ctx context.Context) error

	RunWithGracefulShutdown()
}

// New はサーバーを生成します。
func New(handler http.Handler, opts ...Option) Server {
	server := &server{
		port:            "8080",
		host:            "localhost",
		shutdownTimeout: DefaultShutdownTimeout,
		l:               slog.New(slog.NewTextHandler(os.Stderr, nil)).WithGroup("server"),
	}

	for _, opt := range opts {
		opt(server)
	}

	server.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", server.host, server.port),
		Handler: handler,
	}

	return server
}

// Run はサーバーを起動します。
func (s *server) Run() error {
	s.l.Info(fmt.Sprintf("server starting at %s", s.srv.Addr))
	return s.srv.ListenAndServe()
}

// Shutdown はサーバーを停止します。
func (s *server) Shutdown(ctx context.Context) error {
	s.l.Info("server shutdown ...")
	return s.srv.Shutdown(ctx)
}

// RunWithGracefulShutdown はgraceful shutdownを行うサーバーを起動します。
func (s *server) RunWithGracefulShutdown() {
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
