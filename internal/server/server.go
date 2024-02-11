package server

import (
	"context"
	"time"
)

// DefaultShutdownTimeout はデフォルトのシャットダウンタイムアウトです。
const DefaultShutdownTimeout time.Duration = 10

// Server はサーバーを表すインターフェースです。
type Server interface {
	Run() error
	Shutdown(ctx context.Context) error

	RunWithGracefulShutdown()
}
