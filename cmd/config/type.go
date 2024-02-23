package config

// config はアプリケーションの設定を表す構造体です。基本的には環境変数から読み込みます。
type config struct {
	Application struct {
		ApplicationName string `env:"APPLICATION_NAME" envDefault:"gantrycd"`
		ExternalDomain  string `env:"EXTERNAL_DOMAIN" envDefault:"hogehoge.com"`
	}

	Bff struct {
		Host              string `env:"BFF_HOST" envDefault:"0.0.0.0"`
		Port              int    `env:"BFF_PORT" envDefault:"8080"`
		K8SControllerAddr string `env:"K8S_CONTROLLER_ADDR" envDefault:"localhost:8081"`
		Environment       string `env:"ENVIRONMENTS" envDefault:"local"`
		AllowOrigins      string `env:"ALLOW_ORIGINS" envDefault:"http://localhost:3000"`
	}

	Controller struct {
		Host         string `env:"CONTROLLER_HOST" envDefault:"0.0.0.0"`
		Port         int    `env:"CONTROLLER_PORT" envDefault:"8081"`
		ImageBuilder string `env:"IMAGE_BUILDER" envDefault:"ghcr.io/gantrycd/image-builder:v0.0.3"`
	}

	GitHub struct {
		AppID    int64  `env:"GITHUB_APP_ID" envDefault:"0"`
		CrtPath  string `env:"GITHUB_CRT_PATH" envDefault:""`
		Username string `env:"GITHUB_USERNAME" envDefault:""`
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

	Registry struct {
		Host     string `env:"REGISTRY_HOST" envDefault:"localhost"`
		User     string `env:"REGISTRY_USER" envDefault:""`
		Password string `env:"REGISTRY_PASSWORD" envDefault:""`
	}
}

// Config は読み込まれた設定を保持します。
var Config *config
