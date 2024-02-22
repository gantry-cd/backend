package logger

import (
	"context"
	"fmt"
	"github.com/gantrycd/backend/internal/models"
	v1 "github.com/gantrycd/backend/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			resp, err := logs.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}
				if status.Code(err) == codes.Canceled {
					return
				}
				log.Printf("ストリームからの受信に失敗しました: %v\n", err)
				return
			}

			_, err = fmt.Fprintf(w, "data: %s\n\n", resp.Message)
			if err != nil {
				log.Printf("データの書き込みに失敗しました: %v\n", err)
				return
			}
			flusher.Flush()
		}
	}()
	select {
	case <-ctx.Done():
		if err = logs.CloseSend(); err != nil {
			log.Printf("ログストリームのクローズに失敗しました: %v\n", err)
		}
		<-done
	case <-done:
		break
	}
	return nil
}
