package constants

import "fmt"

// Environment variable lookup strings
const (
	ENV_REST_ADDR    = "REST_ADDR"
	ENV_REDIS_ADDR   = "REDIS_ADDR"
	ENV_GITHUB_TOKEN = "GITHUB_ACCESS_TOKEN"

	ENV_VAULT_ADDR      = "VAULT_ADDR"
	ENV_VAULT_TOKEN     = "VAULT_TOKEN"
	ENV_VAULT_NAMESPACE = "VAULT_NAMESPACE"
)

// EnvNotSetErr returns an error for when an ennvironment variable is s er
func EnvNotSetErr(envVar string) error {
	return fmt.Errorf("%s not set", envVar)
}
