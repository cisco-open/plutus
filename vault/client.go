package vault

import (
	"os"
	"plutus/constants"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// client is an adapter for api.client
type client struct {
	client *api.Client
	logger *logrus.Logger
}

// newClient returns a new Client object
func newClient(namespace string, logger *logrus.Logger) (*client, error) {
	vaultAddr := constants.ENV_VAULT_ADDR
	address, ok := os.LookupEnv(vaultAddr)
	if !ok {
		return nil, errors.Wrap(constants.EnvNotSetErr(vaultAddr), "error creating Vault Client")
	}
	config := &api.Config{
		Address: address,
	}

	vaultClient, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	vaultToken := constants.ENV_VAULT_TOKEN
	token, ok := os.LookupEnv(vaultToken)
	if !ok {
		return nil, errors.Wrap(constants.EnvNotSetErr(vaultToken), "error creating Vault Client")
	}

	vaultClient.SetToken(token)
	vaultClient.SetNamespace(namespace)

	return &client{client: vaultClient, logger: logger}, nil
}
