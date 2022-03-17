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

package cli

import (
	"errors"
	"os"

	"fmt"
	"plutus/config"
	"plutus/constants"
	"plutus/redis"
	"plutus/rest"
	"plutus/vault"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func refreshAction(redisClient *redis.Client, logger *logrus.Logger) error {
	namespaces, err := config.GetNamespaces()
	if err != nil {
		return err
	}

	fmt.Println("here1")

	for _, namespace := range namespaces {
		logger.WithFields(logrus.Fields{
			"namespace": namespace,
			"type":      "refresh",
		}).Info("refreshing data")
		data, err := vault.InitData(namespace, logger)
		if err != nil {
			return err
		}
		redisClient.RefreshAll(data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) refreshAction(*cli.Context) error {
	return refreshAction(a.redisClient, a.logger)
}

func (a *App) startAPIServer(c *cli.Context) error {

	server, err := rest.NewAPIServer(a.redisClient, a.groupsReader, a.logger)
	if err != nil {
		return err
	}

	err = refreshAction(a.redisClient, a.logger)
	if err != nil {
		return err
	}
	server.Run(refreshAction)
	return nil
}

func getShellNamespace() (string, error) {
	vaultNamespace := constants.ENV_VAULT_NAMESPACE
	namespace, ok := os.LookupEnv(vaultNamespace)
	if !ok {
		return "", constants.EnvNotSetErr(vaultNamespace)
	}
	return namespace, nil
}

func (a *App) lookupPoliciesAction(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("no Username provided")
	}

	username := c.Args().First()
	namespace, err := getShellNamespace()
	if err != nil {
		return err
	}

	pols, err := a.redisClient.QueryPoliciesFromUsername(namespace, username)
	if err != nil {
		return err
	}
	fmt.Println(pols)
	return nil
}

func (a *App) lookupUsersAction(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("no vault path prided")
	}

	namespace, err := getShellNamespace()
	if err != nil {
		return err
	}
	firstArg := c.Args().First()
	var users []string
	if c.Bool("policy") {
		users, err = a.redisClient.QueryUsersFromPolicyName(namespace, firstArg)

		if err != nil {
			return err
		}
	} else if c.Bool("path") {
		usersSet, err := a.redisClient.QueryUsersFromVaultPath(namespace, firstArg)
		users = usersSet.Users()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("please specify if the provided argument \"%s\" is a vault path or a policy name using the flags", firstArg)
	}
	fmt.Println(users)
	return nil
}
