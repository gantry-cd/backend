package config

import "time"

// config はアプリケーションの設定を表す構造体です。基本的には環境変数から読み込みます。
type config struct {
	Server struct {
		Host            string        `env:"HOST" envDefault:"localhost"`
		Port            int           `env:"PORT" envDefault:"8080"`
		ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
	}

	Application struct {
		Name string `env:"APP_NAME" envDefault:"gantry"`
	}

	KeyCloak struct {
		KeyCloakHost  string `env:"KEYCLOAK_HOST" envDefault:"http://localhost:8080"`
		Realm         string `env:"KEYCLOAK_REALM" envDefault:"gantry"`
		ClientID      string `env:"KEYCLOAK_CLIENT_ID" envDefault:"gantry-app"`
		RedirectURIs  string `env:"KEYCLOAK_REDIRECT_URIS" envDefault:"/api/auth/callback/keycloak"`
		WebOrigins    string `env:"KEYCLOAK_WEB_ORIGINS" envDefault:"http://localhost:3000"`
		RootURL       string `env:"KEYCLOAK_ROOT_URL" envDefault:"http://localhost:3000"`
		AdminID       string `env:"KEYCLOAK_ADMIN_ID" envDefault:"admin"`
		AdminPassword string `env:"KEYCLOAK_ADMIN_PASSWORD" envDefault:"admin"`
	}
}

// Config は読み込まれた設定を保持します。
var Config *config
