package github

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gantrycd/backend/internal/usecases/application/controller"
	"github.com/google/go-github/v29/github"
)

type handler struct {
	l *slog.Logger

	interactor controller.GithubAppEvents
}

// Option はサーバーのオプションを設定するための関数です。
type Option func(*handler)

// WithLogger はロガーを設定するオプションです。
func WithLogger(l *slog.Logger) Option {
	return func(s *handler) {
		s.l = l
	}
}

// New は新しいWebhookのハンドラーを作成します。
func New(interactor controller.GithubAppEvents, opts ...Option) *handler {
	handler := &handler{
		l:          slog.New(slog.NewTextHandler(os.Stderr, nil)).WithGroup("handler"),
		interactor: interactor,
	}

	for _, opt := range opts {
		opt(handler)
	}

	return handler
}

func (h *handler) GithubAppsHandler(w http.ResponseWriter, r *http.Request) {
	event, err := h.parseWebhook(r)
	if err != nil {
		h.l.Error("error parsing webhook", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ここでeventを処理する
	switch e := event.(type) {
	case *github.InstallationEvent:
		err = h.Installation(e)
	case *github.InstallationRepositoriesEvent:
		err = h.InstallationRepositories(e)
	case *github.MetaEvent:
		err = h.Meta(e)
	case *github.CreateEvent:
		err = h.Create(e)
	case *github.DeleteEvent:
		err = h.Delete(e)
	case *github.PushEvent:
		err = h.Push(e)
	case *github.PullRequestEvent:
		err = h.PullRequest(e)
	case *github.RepositoryEvent:
		err = h.Repository(e)
	default:
		h.l.Info(fmt.Sprintf("event not supported: %v", event))
		http.Error(w, "event not supported", http.StatusBadRequest)
	}

	if err != nil {
		h.l.Error("error processing event", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// parseWebhook はWebhookのペイロードをパースして任意のGithubイベントに変換します。
func (h *handler) parseWebhook(r *http.Request) (any, error) {
	payload, err := github.ValidatePayload(r, nil)
	if err != nil {
		return nil, err
	}
	return github.ParseWebHook(github.WebHookType(r), payload)
}
