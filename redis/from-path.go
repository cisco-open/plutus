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

import (
	"plutus/parser"
	"plutus/utils"
	"strings"
)

// QueryUsersFromVaultPath returns the list of users that can access the given vault path
// Given a path "a/b/c"
// It will query the paths ["a/*", "a/b/*", "a/b/c"]
func (c *Client) QueryUsersFromVaultPath(namespace, path string) (UserResponseSet, error) {
	usersSet, err := c.queryUsersFromVaultPath(namespace, path)
	if err != nil {
		return nil, err
	}

	slash := "/"
	splitPath := strings.Split(path, slash)

	for i := 0; i < len(splitPath); i++ {
		upperPath := strings.Join(splitPath[:i], slash) + "/*"
		upperPathSet, err := c.queryUsersFromVaultPath(namespace, upperPath)
		if err != nil {
			return nil, err
		}
		usersSet.AddSet(upperPathSet)
	}

	return usersSet, nil
}

func (c *Client) queryUsersFromVaultPath(namespace, path string) (UserResponseSet, error) {
	usersSet := make(UserResponseSet, 0)

	encodedPolicies, err := c.smembers(namespace, PrefixPatToPol, path)
	if err != nil {
		return nil, err
	}

	decodedPaths, err := parser.DecodePaths(encodedPolicies)
	if err != nil {
		return nil, err
	}

	for _, decoded := range decodedPaths {
		policyName := decoded.PolicyName

		decodedCapabilities := decoded.Capabilities

		policyUsers, err := c.QueryUsersFromPolicyName(namespace, policyName)
		if err != nil {
			return nil, err
		}

		for _, policyUser := range policyUsers {
			pr, err := c.getPolicyUserRelation(namespace, policyUser, policyName)
			pr.Capabilities = decodedCapabilities.Strings()
			if err != nil {
				return nil, err
			}
			usersSet.AddPolicyResponse(policyUser, pr)
		}
	}
	return usersSet, nil
}

func (c *Client) getPolicyUserRelation(namespace, username, policyName string) (PolicyResponse, error) {
	pr := NewPolicyResponse(policyName)

	// CHECKING VAULT GROUPS
	polGroups, err := c.smembers(namespace, PrefixPolToVgr, policyName)
	if err != nil {
		return PolicyResponse{}, err
	}
	userGroups, err := c.smembers(namespace, PrefixUsrToVgr, username)
	if err != nil {
		return PolicyResponse{}, err
	}
	pr.VaultGroups = utils.Intersection(polGroups, userGroups)

	// CHECKING VAULT ENTITIES
	userEntity, err := c.get(namespace, PrefixUsrToEnt, username)
	if err == nil {
		polEntities, err := c.smembers(namespace, PrefixPolToEnt, policyName)
		if err != nil {
			return PolicyResponse{}, err
		}
		for _, entity := range polEntities {
			if entity == userEntity {
				pr.Entities = append(pr.Entities, entity)
			}
		}
	}

	// CHECKING OIDC ROLES
	polRoles, err := c.smembers(namespace, PrefixPolToRol, policyName)
	if err != nil {
		return PolicyResponse{}, err
	}
	userRoles, err := c.QueryRolesFromUsername(namespace, username)
	if err != nil {
		return PolicyResponse{}, err
	}
	pr.OIDCRoles = utils.Intersection(userRoles, polRoles)

	return pr, nil
}
