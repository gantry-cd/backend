package webhook

import (
	"log/slog"
	"net/http"

	"github.com/google/go-github/v29/github"
)

type handler struct {
	mux *http.ServeMux
	l   *slog.Logger
}

// New は新しいWebhookのハンドラーを作成します。
func New() http.Handler {
	handler := &handler{
		mux: http.NewServeMux(),
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

	h.l.Debug("event received", event)

	// ここでeventを処理する
	switch e := event.(type) {

	default:
		h.l.Info("event not supported", "event", e)
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
