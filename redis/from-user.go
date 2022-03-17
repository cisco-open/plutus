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
)

// QueryGroupsFromUsername returns the groups the given username is a part of
func (c *Client) QueryGroupsFromUsername(namespace, username string) ([]string, error) {
	return c.smembers(namespace, PrefixUsrToVgr, username)
}

// QueryRolesFromUsername returns the roles the given username has
func (c *Client) QueryRolesFromUsername(namespace, username string) ([]string, error) {
	externalGroups, err := c.smembers(namespace, PrefixUsrToEgr, username)
	if err != nil {
		return nil, err
	}

	roles := make([]string, 0)
	for _, egr := range externalGroups {
		_roles, err := c.smembers(namespace, PrefixEgrToRol, egr)
		if err != nil {
			return nil, err
		}
		roles = append(roles, _roles...)
	}

	return roles, nil
}

// QueryPathsFromUsername returns the paths the given username can access
func (c *Client) QueryPathsFromUsername(namespace, username string) (PathResponseSet, error) {

	pols, err := c.QueryPoliciesFromUsername(namespace, username)
	if err != nil {
		return nil, err
	}

	pathsSet := make(PathResponseSet, 0)
	for _, pol := range pols.AsSlice() {
		encodedPolPaths, err := c.smembers(namespace, PrefixPolToPat, pol.Name)
		if err != nil {
			return nil, err
		}

		for _, encPath := range encodedPolPaths {

			path, err := parser.DecodePath(encPath)
			if err != nil {
				return nil, err
			}
			pathsSet[path.Path] = NewPathResponse(path)
		}
	}

	return pathsSet, nil
}

// QueryPoliciesFromUsername returns the policies the given username has
func (c *Client) QueryPoliciesFromUsername(namespace, username string) (PolicyResponseSet, error) {
	policiesSet := make(PolicyResponseSet)

	groups, err := c.QueryGroupsFromUsername(namespace, username)
	if err == nil {

		for _, group := range groups {
			grpPolicies, err := c.smembers(namespace, PrefixVgrToPol, group)
			if err != nil {
				return PolicyResponseSet{}, err
			}
			for _, grpPolicy := range grpPolicies {
				policiesSet.Add(PolicyResponse{
					Name:        grpPolicy,
					VaultGroups: []string{group},
				})
			}
		}
	}

	entityID, err := c.get(namespace, PrefixUsrToEnt, username)
	if err == nil {
		entityPolicies, err := c.smembers(namespace, PrefixEntToPol, entityID)
		if err != nil {
			return PolicyResponseSet{}, err
		}
		for _, entityPolicy := range entityPolicies {
			policiesSet.Add(PolicyResponse{
				Name:     entityPolicy,
				Entities: []string{entityID},
			})
		}
	}

	roles, err := c.QueryRolesFromUsername(namespace, username)
	if err == nil {
		for _, role := range roles {
			rolePolicies, err := c.smembers(namespace, PrefixRolToPol, role)
			if err != nil {
				return PolicyResponseSet{}, err
			}
			for _, rolePolicy := range rolePolicies {
				policiesSet.Add(PolicyResponse{
					Name:      rolePolicy,
					OIDCRoles: []string{role},
				})
			}
		}
	}

	return policiesSet, nil
}
