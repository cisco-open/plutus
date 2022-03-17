// Copyright 2022 Cisco Systems, Inc. and its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

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
