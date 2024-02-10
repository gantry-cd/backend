package githubapp

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	customController "github.com/gantrycd/backend/proto/k8s-controller"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DefaultNamespacePrefix = "gantrycd"
)

type githubAppEvents struct {
	l *slog.Logger

	customController customController.K8SCustomControllerClient
}

// githubAppEvents はGithubAppのインタラクターのインターフェースです。
type GithubAppEvents interface {
	CreateNameSpace(ctx context.Context, organization string) error
	ListNameSpace(ctx context.Context, prefix string) ([]string, error)
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
func New(customController customController.K8SCustomControllerClient, opts ...Option) GithubAppEvents {
	ge := &githubAppEvents{
		customController: customController,
		l:                slog.New(slog.NewTextHandler(os.Stderr, nil)).WithGroup("app-interactor"),
	}

	for _, opt := range opts {
		opt(ge)
	}

	return ge
}

// CreateNameSpace はOrganization名を元にNamespaceを作成します。
func (ge *githubAppEvents) CreateNameSpace(ctx context.Context, organization string) error {
	// GitHubのOrganization名をRFC1123に基づいてNamespace名として利用する
	// GitHubのOrganization名は大文字小文字が混在しているため、小文字に変換する。あとは基本的にRFC1123に基づいているため、そのまま利用する。
	orgName := strings.ToLower(organization)

	_, err := ge.customController.CreateNamespace(ctx, &customController.CreateNamespaceRequest{
		Name: fmt.Sprintf("%s-%s", DefaultNamespacePrefix, orgName),
	})

	status, _ := status.FromError(err)

	// Namespaceが既に存在する場合はエラーを無視する
	if err != nil || status.Code() != codes.AlreadyExists {
		return err
	}

	return nil
}

// ListNameSpace はNamespaceの一覧を取得します。
func (ge *githubAppEvents) ListNameSpace(ctx context.Context, prefix string) ([]string, error) {
	result, err := ge.customController.ListNamespaces(ctx, nil)
	if err != nil {
		ge.l.Error("error listing namespaces", "error", err.Error())
		return nil, err
	}

	var namespaces []string
	for _, ns := range result.Names {
		if strings.HasPrefix(ns, prefix) {
			namespaces = append(namespaces, ns)
		}
	}

	return namespaces, nil
}
