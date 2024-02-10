package ghapp

import (
	"log/slog"
	"os"

	"github.com/google/go-github/v29/github"
)

type githubAppInteractor struct {
	l *slog.Logger
}

// GithubAppInteractor はGithubAppのインタラクターのインターフェースです。
type GithubAppInteractor interface {
	InstallApplication(*github.InstallationEvent) error
	UninstallApplication(*github.InstallationEvent) error

	RepositoriesAdded(*github.InstallationRepositoriesEvent) error
	RepositoriesRemoved(*github.InstallationRepositoriesEvent) error

	CreatePullRequest(*github.CreateEvent) error
	DeletePullRequest(*github.DeleteEvent) error

	Push(*github.PushEvent) error
}

// Option はサーバーのオプションを設定するための関数です。
type Option func(*githubAppInteractor)

// WithLogger はロガーを設定するオプションです。
func WithLogger(l *slog.Logger) Option {
	return func(s *githubAppInteractor) {
		s.l = l
	}
}

// New は新しいGithubAppのインタラクターを作成します。
func New(opts ...Option) GithubAppInteractor {
	gi := &githubAppInteractor{
		l: slog.New(slog.NewTextHandler(os.Stderr, nil)).WithGroup("app-interactor"),
	}

	for _, opt := range opts {
		opt(gi)
	}

	return gi
}

func (gi *githubAppInteractor) InstallApplication(e *github.InstallationEvent) error {
	panic("not implemented") // TODO: Implement
}

func (gi *githubAppInteractor) UninstallApplication(*github.InstallationEvent) error {
	panic("not implemented") // TODO: Implement
}

func (gi *githubAppInteractor) RepositoriesAdded(*github.InstallationRepositoriesEvent) error {
	panic("not implemented") // TODO: Implement
}

func (gi *githubAppInteractor) RepositoriesRemoved(*github.InstallationRepositoriesEvent) error {
	panic("not implemented") // TODO: Implement
}

func (gi *githubAppInteractor) CreatePullRequest(*github.CreateEvent) error {
	panic("not implemented") // TODO: Implement
}

func (gi *githubAppInteractor) DeletePullRequest(*github.DeleteEvent) error {
	panic("not implemented") // TODO: Implement
}

func (gi *githubAppInteractor) Push(*github.PushEvent) error {
	panic("not implemented") // TODO: Implement
}
