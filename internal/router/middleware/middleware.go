package middleware

import (
	"log/slog"
	"net/http"
	"os"
)

type middleware struct {
	l *slog.Logger
}

type Middleware interface {
	// KeyCloakOAuth はKeyCloakのOAuth認証を行うミドルウェアを返す
	KeyCloakOAuth(h http.HandlerFunc) http.HandlerFunc
	// Integrate は環境に応じたミドルウェアを返す
	Integrate(h http.HandlerFunc) http.Handler
	// Recover はpanicをrecoverするミドルウェアを返す
	Recover(h http.Handler) http.Handler
	// AccessLogger はアクセスログを出力するミドルウェアを返す
	AccessLogger(h http.Handler) http.Handler
	// AllowAllOrigins は全てのオリジンからのリクエストを許可するミドルウェアを返す
	AllowAllOrigins(h http.Handler) http.Handler
	// CorsWithEnv は環境に応じたCORS設定を行うミドルウェアを返す
	CorsWithEnv(h http.Handler) http.Handler
}

func NewMiddleware() Middleware {
	return &middleware{
		l: slog.New(slog.NewTextHandler(os.Stdout, nil)).WithGroup("middleware"),
	}
}

func BuildChain(h http.Handler, m ...func(http.Handler) http.Handler) http.Handler {
	if len(m) == 0 {
		return h
	}

	return m[0](BuildChain(h, m[1:cap(m)]...))
}
