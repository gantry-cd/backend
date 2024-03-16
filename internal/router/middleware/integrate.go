package middleware

import (
	"net/http"

	"github.com/aura-cd/backend/cmd/config"
)

func (m *middleware) Integrate(h http.HandlerFunc) http.Handler {
	switch config.Config.Bff.Environment {
	case "local":
		return BuildChain(h, m.Recover, m.AccessLogger, m.AllowAllOrigins)
	case "dev":
		return BuildChain(h, m.Recover, m.AccessLogger, m.CorsWithEnv)
	case "staging":
		return BuildChain(h, m.Recover, m.AccessLogger, m.CorsWithEnv)
	case "production":
		return BuildChain(h, m.Recover, m.AccessLogger, m.CorsWithEnv)
	default:
		return BuildChain(h, m.Recover, m.AccessLogger, m.AllowAllOrigins)
	}
}
