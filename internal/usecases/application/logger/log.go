package logger

import (
	"context"
	"fmt"
	"github.com/gantrycd/backend/internal/models"
	v1 "github.com/gantrycd/backend/proto"
	"io"
	"log"
	"net/http"
)

type logInteractor struct {
	resource v1.K8SCustomControllerClient
}
type LogInteractor interface {
	GetLogStream(ctx context.Context, w http.ResponseWriter, request models.PodLogRequest) error
}

func New(resource v1.K8SCustomControllerClient) LogInteractor {
	return &logInteractor{
		resource: resource,
	}
}

func (r *logInteractor) GetLogStream(ctx context.Context, w http.ResponseWriter, request models.PodLogRequest) error {
	flusher, _ := w.(http.Flusher)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	logs, err := r.resource.GetLogs(ctx, &v1.GetLogsRequest{
		Namespace: request.Organization,
		PodName:   request.Pod,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	go func() {
		for {
			resp, err := logs.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("ストリームからの受信に失敗しました: %v", err)
			}

			_, err = fmt.Fprintf(w, "data: %s\n\n", resp.Message)
			if err != nil {
				return
			}
			flusher.Flush()
		}
	}()
	<-ctx.Done()
	return nil
}
