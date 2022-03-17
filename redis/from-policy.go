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

package redis

// QueryUsersFromPolicyName gets the Users that have the provided policy
func (c *Client) QueryUsersFromPolicyName(namespace, policyName string) ([]string, error) {
	usersSet := make(map[string]bool, 0)

	groups, err := c.smembers(namespace, PrefixPolToVgr, policyName)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		grpUsers, err := c.smembers(namespace, PrefixVgrToUsr, group)
		if err != nil {
			return nil, err
		}

		for _, grpUser := range grpUsers {
			usersSet[grpUser] = true
		}
	}

	entities, err := c.smembers(namespace, PrefixPolToEnt, policyName)
	if err != nil {
		return nil, err
	}

	for _, entityID := range entities {
		username, err := c.get(namespace, PrefixEntToUsr, entityID)
		if err != nil {
			return nil, err
		}
		usersSet[username] = true

	}

	roles, err := c.smembers(namespace, PrefixPolToRol, policyName)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		externalGroups, err := c.smembers(namespace, PrefixRolToEgr, role)
		if err != nil {
			return nil, err
		}

		for _, externalGroup := range externalGroups {
			externalUsers, err := c.smembers(namespace, PrefixEgrToUsr, externalGroup)
			if err != nil {
				return nil, err
			}

			for _, externalUser := range externalUsers {
				usersSet[externalUser] = true
			}
		}
	}

	i := 0
	usersSlice := make([]string, len(usersSet))
	for username := range usersSet {
		usersSlice[i] = username
		i++
	}

	return usersSlice, nil
}
