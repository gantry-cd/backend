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
