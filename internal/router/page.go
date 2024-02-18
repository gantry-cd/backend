package router

import "net/http"

func (r *router) page() {
	r.mux.Handle("GET /", r.middleware.KeyCloakOAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	})))
}
