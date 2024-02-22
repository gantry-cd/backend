package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v13"
)

// envに起こそう
var (
	keycloakDomain       = "http://10.10.10.40:8080"
	keycloakRealm        = "test"
	keycloakClientID     = "test-client"
	keycloakClientSecret = "LntMHtmL5Ef2KXtIy0u1TNG8AJBgRVW0"
)

const (
	AuthorizationHeaderKey = "Authorization"
	AuthorizationType      = "Bearer"
)

func (m *middleware) KeyCloakOAuth(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := gocloak.NewClient(keycloakDomain)
		ctx := context.Background()
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

		rptResult, err := client.RetrospectToken(ctx, auth, keycloakClientID, keycloakClientSecret, keycloakRealm)
		if err != nil {
			log.Println("Failed to introspect token", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !*rptResult.Active {
			log.Println("Token is not active")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// role, err := client.GetRoleMappingByUserID(ctx, token.AccessToken, keycloakRealm, keycloakClientID)
		// if err != nil {
		// 	log.Println("Failed to get role mapping", err)
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	return
		// }
	})

}
