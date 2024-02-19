package config

// config はアプリケーションの設定を表す構造体です。基本的には環境変数から読み込みます。
type config struct {
	Bff struct {
		Host              string `env:"BFF_HOST" envDefault:"0.0.0.0"`
		Port              int    `env:"BFF_PORT" envDefault:"8080"`
		K8SControllerAddr string `env:"K8S_CONTROLLER_ADDR" envDefault:"localhost:8081"`
	}

	Controller struct {
		Host string `env:"CONTROLLER_HOST" envDefault:"0.0.0.0"`
		Port int    `env:"CONTROLLER_PORT" envDefault:"8081"`
	}

	GitHub struct {
		AppID     int64  `env:"GITHUB_APP_ID" envDefault:"0"`
		InstallID int64  `env:"GITHUB_INSTALL_ID" envDefault:"0"`
		CrtPath   string `env:"GITHUB_CRT_PATH" envDefault:""`
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
