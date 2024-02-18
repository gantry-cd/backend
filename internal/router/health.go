package router

import "net/http"

func (r *router) health() {
	r.mux.Handle("GET /health", (http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	})))
}
