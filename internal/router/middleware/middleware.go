package middleware


import (
	"net/http"
)

type middleware struct {
}

type Middleware interface {
	KeyCloakOAuth(h http.HandlerFunc) http.HandlerFunc
}

func NewMiddleware() Middleware {
	return &middleware{}
}

func BuildChain(h http.Handler, m ...func(http.Handler) http.Handler) http.Handler {
	if len(m) == 0 {
		return h
	}

	return m[0](BuildChain(h, m[1:cap(m)]...))
}
