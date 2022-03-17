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
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Data contains all the data from Vault for a given namespace
// Each data struct has it's own vault client as well
type Data struct {
	client    *client           // vault client
	Namespace string            // namespace
	Groups    map[string]Group  // group name to group
	Entities  map[string]Entity // entityID to entity
	Policies  map[string]Policy // policy name to policy
	Aliases   map[string]Alias  // entityID to alias
	Roles     map[string]Role   // role name to role
}

// InitData initialises data with a vault client
// if fetchAll is true, it fetches the data as well
func InitData(namespace string, logger *logrus.Logger) (*Data, error) {
	vaultClient, err := newClient(namespace, logger)
	if err != nil {
		return nil, err
	}

	data := &Data{
		client:    vaultClient,
		Namespace: namespace,
	}

	err = data.FetchAllData(logger)
	if err != nil {
		return nil, errors.Wrapf(err, "err fetching data for namespace %s", namespace)
	}

	return data, nil
}

// FetchAllData fetches all the data by fetching it from vault
// Note that this only stores the data in the Data struct and does not update the redis instance
// To update the redis instance, use the Refresh* receivers on the redisClient struct
func (d *Data) FetchAllData(logger *logrus.Logger) error {
	logger.WithFields(logrus.Fields{
		"namespace": d.Namespace,
		"type":      "fetch",
	}).Info("Fetching all data")
	client := d.client

	logger.Info("fetching policies")
	policies, err := client.Policies()
	if err != nil {
		return err
	}
	d.Policies = policies
	logger.Info("Number of policies fetched: ", len(policies))

	logger.Info("fetching groups")
	groups, err := client.Groups()
	if err != nil {
		return err
	}
	d.Groups = groups
	logger.Info("Number of groups fetched: ", len(groups))

	logger.Info("fetching entities")
	entities, err := client.Entities()
	if err != nil {
		return err
	}
	d.Entities = entities
	logger.Info("Number of entities fetched: ", len(entities))

	logger.Info("fetching aliases")
	aliases, err := client.Aliases()
	if err != nil {
		return err
	}
	d.Aliases = aliases
	logger.Info("Number of aliases fetched: ", len(aliases))

	logger.Info("fetching roles")
	roles, err := client.Roles()
	if err != nil {
		return err
	}
	d.Roles = roles
	logger.Info("Number of roles fetched: ", len(roles))

	return nil

}
