package webhook

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/go-github/v29/github"
)

type handler struct {
	mux *http.ServeMux
	l   *slog.Logger
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
func New(opts ...Option) http.Handler {
	handler := &handler{
		mux: http.NewServeMux(),
		l:   slog.New(slog.NewTextHandler(os.Stderr, nil)).WithGroup("handler"),
	}

	for _, opt := range opts {
		opt(handler)
	}

	handler.mux.Handle("POST /github/app/webhook", http.HandlerFunc(handler.githubAppsHandler))
	return handler.mux
}

func (h *handler) githubAppsHandler(w http.ResponseWriter, r *http.Request) {
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
