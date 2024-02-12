package middleware

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// envに起こそう
var (
	keycloakDomain       = "your-keycloak-domain"
	keycloakRealm        = "your-keycloak-realm"
	keycloakClientID     = "your-client-id"
	keycloakClientSecret = "your-client-secret"
	keycloakRequestURL   = fmt.Sprintf("http://%s/auth/realms/%s/protocol/openid-connect/token/introspect", keycloakDomain, keycloakRealm)
)

const (
	AuthorizationHeaderKey = "Authorization"
	AuthorizationType      = "Bearer"
)

func (m *middleware) KeyCloakOAuth() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get(AuthorizationHeaderKey)
		if auth == "" {
			log.Println("Authorization header is not found")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authParts := strings.Split(auth, " ")
		if len(authParts) != 2 {
			log.Println("Authorization header is invalid")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if authParts[0] != AuthorizationType {
			log.Println("Authorization type is invalid")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		params := url.Values{}
		token := authParts[1]
		params.Set("token", token)
		params.Set("token_hint", "access_token")
		params.Set("client_id", keycloakClientID)
		params.Set("client_secret", keycloakClientSecret)
		_, err := http.NewRequest("POST", keycloakRequestURL, nil)
		if err != nil {
			log.Println("Failed to create request to Keycloak", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

	})

}
