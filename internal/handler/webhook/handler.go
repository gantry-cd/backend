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
		h.l.Info(fmt.Sprintf("installation event received: %v", e))
	case *github.InstallationRepositoriesEvent:
		h.l.Info(fmt.Sprintf("installation repositories event received: %v", e))
	case *github.MetaEvent:
		h.l.Info(fmt.Sprintf("meta event received: %v", e))
	case *github.CreateEvent:
		h.l.Info(fmt.Sprintf("create event received: %v", e))
	case *github.DeleteEvent:
		h.l.Info(fmt.Sprintf("delete event received: %v", e))
	case *github.PushEvent:
		h.l.Info(fmt.Sprintf("push event received: %v", e))
	case *github.PullRequestEvent:
		h.l.Info(fmt.Sprintf("pull request event received: %v", e))
	case *github.RepositoryEvent:
		h.l.Info(fmt.Sprintf("repository event received: %v", e))
	default:
		h.l.Info(fmt.Sprintf("event not supported: %v", event))
		http.Error(w, "event not supported", http.StatusBadRequest)
	}
}

// parseWebhook はWebhookのペイロードをパースして任意のGithubイベントに変換します。
func (h *handler) parseWebhook(r *http.Request) (any, error) {
	payload, err := github.ValidatePayload(r, nil)
	if err != nil {
		return nil, err
	}
	return github.ParseWebHook(github.WebHookType(r), payload)
}
