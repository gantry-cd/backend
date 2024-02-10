package ghevent

import (
	"log/slog"
	"os"
)

type githubAppEvents struct {
	l *slog.Logger
}

// githubAppEvents はGithubAppのインタラクターのインターフェースです。
type GithubAppEvents interface {
}

// Option はサーバーのオプションを設定するための関数です。
type Option func(*githubAppEvents)

// WithLogger はロガーを設定するオプションです。
func WithLogger(l *slog.Logger) Option {
	return func(s *githubAppEvents) {
		s.l = l
	}
}

// New は新しいGithubAppのインタラクターを作成します。
func New(opts ...Option) GithubAppEvents {
	ge := &githubAppEvents{
		l: slog.New(slog.NewTextHandler(os.Stderr, nil)).WithGroup("app-interactor"),
	}

	for _, opt := range opts {
		opt(ge)
	}

	return ge
}
