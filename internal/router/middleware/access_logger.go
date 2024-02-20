package middleware

import (
	"net/http"
	"time"
)

type accessLog struct {
	Time       time.Time
	RemoteAddr string
	RequestURI string
	Method     string
	Duration   time.Duration
	UserAgent  string
}

func (m *middleware) AccessLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var al accessLog = accessLog{
			Time:       time.Now(),
			RemoteAddr: r.RemoteAddr,
			RequestURI: r.RequestURI,
			Method:     r.Method,
			UserAgent:  r.UserAgent(),
		}
		defer func() {
			al.Duration = time.Since(al.Time)
			m.l.Info("Access Log", "time", al.Time, "remote_addr", al.RemoteAddr, "request_uri", al.RequestURI, "method", al.Method, "duration", al.Duration, "user_agent", al.UserAgent)
		}()

		h.ServeHTTP(w, r)
	})
}
