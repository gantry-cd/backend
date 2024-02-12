package middleware

import "net/http"

type middleware struct {
}

type Middleware interface {
	KeyCloakOAuth() http.HandlerFunc
}

func NewMiddleware() Middleware {
	return &middleware{}
}
