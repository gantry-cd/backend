package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
)

const (
	RealmMaster string = "master"
)

type kcClient struct {
	host string

	client *gocloak.GoCloak
	token  string
}

func (c *kcClient) Token() string {
	return c.token
}

func (c *kcClient) Client() *gocloak.GoCloak {
	return c.client
}

// New は Keycloak クライアントを生成します。
func New(host string) *kcClient {
	client := gocloak.NewClient(host)
	return &kcClient{
		host:   host,
		client: client,
	}
}

// LoginAdmin は Keycloak の master realm の admin としてログインします。
func (c *kcClient) LoginAdmin(ctx context.Context, user, password string) error {
	token, err := c.client.LoginAdmin(context.Background(), user, password, RealmMaster)
	if err != nil {
		return err
	}

	c.token = token.AccessToken
	return nil
}
