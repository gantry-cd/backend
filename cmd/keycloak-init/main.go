package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/aura-cd/backend/cmd/config"
	"github.com/aura-cd/backend/internal/driver/keycloak"
	"github.com/aura-cd/backend/internal/utils/random"
)

func init() {
	config.LoadEnv()
}

const (
	RealmMaster     string = "master"
	DefaultRealm    string = "gantry"
	DefaultProtocol string = "openid-connect"
	DefaultAuthType string = "client-secret"
)

const (
	EnvFileTmlp = `
# Keycloak
KEYCLOAK_HOST=%s
KEYCLOAK_REALM=%s
KEYCLOAK_CLIENT_ID=%s
KEYCLOAK_CLIENT_SECRET=%s
KEYCLOAK_ADMIN_ID=%s
KEYCLOAK_ADMIN_PASSWORD=%s
KEYCLOAK_REDIRECT_URIS=%s
KEYCLOAK_WEB_ORIGIN=%s
KEYCLOAK_ROOT_URL=%s
`
)

func main() {
	kc := keycloak.New(config.Config.KeyCloak.KeyCloakHost)

	ctx := context.Background()

	if err := kc.LoginAdmin(ctx, config.Config.KeyCloak.AdminID, config.Config.KeyCloak.AdminPassword); err != nil {
		panic(err)
	}

	client := kc.Client()

	realm, err := client.GetRealm(ctx, kc.Token(), DefaultRealm)
	if err != nil {
		if apiErr, ok := err.(*gocloak.APIError); !ok || apiErr.Code != 404 {
			panic(err)
		}
	}

	if realm != nil {
		fmt.Println("already exists")
		os.Exit(1)
	}

	// Realmの作成
	if _, err = client.CreateRealm(ctx, kc.Token(), gocloak.RealmRepresentation{
		Realm:   toPtr(DefaultRealm),
		Enabled: toPtr(true),
	}); err != nil {
		panic(err)
	}

	// Clientの作成
	clientSecret, err := random.RandomString(20)
	if err != nil {
		panic(err)
	}

	clientID, err := client.CreateClient(ctx, kc.Token(), DefaultRealm, gocloak.Client{
		ClientID:                     toPtr(config.Config.KeyCloak.ClientID),
		Secret:                       toPtr(clientSecret),
		Name:                         toPtr(config.Config.KeyCloak.ClientID),
		Description:                  toPtr("GantryCD内の認証認可を行うためのクライアント"),
		PublicClient:                 toPtr(false),
		AuthorizationServicesEnabled: toPtr(true),
		ServiceAccountsEnabled:       toPtr(true),
		StandardFlowEnabled:          toPtr(true),
		ClientAuthenticatorType:      toPtr(DefaultAuthType),
		Protocol:                     toPtr(DefaultProtocol),
		RedirectURIs:                 toPtr(strings.Split(config.Config.KeyCloak.RedirectURIs, ",")),
		WebOrigins:                   toPtr(strings.Split(config.Config.KeyCloak.WebOrigins, ",")),
		RootURL:                      toPtr(config.Config.KeyCloak.RootURL),
		Enabled:                      toPtr(true),
	})

	if err != nil {
		panic(err)
	}

	// Adminの作成
	adminID, err := client.CreateUser(ctx, kc.Token(), DefaultRealm, gocloak.User{
		Username: toPtr(config.Config.KeyCloak.AdminID),
		ID:       toPtr(config.Config.KeyCloak.AdminID),
		Enabled:  toPtr(true),
	})

	if err != nil {
		panic(err)
	}

	// Adminのパスワード設定
	adminPassword, err := random.RandomString(20)
	if err != nil {
		panic(err)
	}

	// パスワードを変更させるように設定しておく
	if err := client.SetPassword(ctx, kc.Token(), adminID, DefaultRealm, adminPassword, true); err != nil {
		panic(err)
	}

	// 設定に応じたenvファイルを作成
	envFile, err := os.Create("keycloak.env")
	if err != nil {
		panic(err)
	}

	defer envFile.Close()

	_, err = envFile.WriteString(
		fmt.Sprintf(
			EnvFileTmlp,
			config.Config.KeyCloak.KeyCloakHost,
			DefaultRealm,
			clientID,
			clientSecret,
			config.Config.KeyCloak.AdminID,
			adminPassword,
			config.Config.KeyCloak.RedirectURIs,
			config.Config.KeyCloak.WebOrigins,
			config.Config.KeyCloak.RootURL,
		),
	)

	if err != nil {
		panic(err)
	}

}

func toPtr[T any](i T) *T {
	return &i
}
