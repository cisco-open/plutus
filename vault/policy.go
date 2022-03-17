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
	"fmt"
	"plutus/parser"

	"github.com/pkg/errors"
)

// Policy is used to represent a vault policy
type Policy struct {
	Name    string
	Content string
	Paths   []parser.ParsedPath
}

// EncodedPaths returns the list of paths encoded as strings
// ready to be added to the redis instance
func (p *Policy) EncodedPaths() ([]string, error) {
	encodedPaths := make([]string, len(p.Paths))

	for i, path := range p.Paths {
		encodedPath, err := path.Encode()
		if err != nil {
			return nil, err
		}
		encodedPaths[i] = encodedPath
	}
	return encodedPaths, nil
}

// Policies returns a map of all policies on vault
// The mapping is from the policy name to the policy object
func (c *client) Policies() (map[string]Policy, error) {
	sys := c.client.Sys()

	policyNames, err := sys.ListPolicies()
	if err != nil {
		return nil, err
	}

	policies := make(map[string]Policy)

	for _, policyName := range policyNames {
		content, err := sys.GetPolicy(policyName)
		if err != nil {
			return nil, err
		}
		policies[policyName], err = newPolicy(policyName, content)
		if err != nil {
			c.logger.Warning(err)
		}
	}

	return policies, nil

}

func newPolicy(name, content string) (Policy, error) {
	hclParser := parser.NewHCLParser()
	paths, err := hclParser.ParseString(content, name)
	if err != nil {
		return Policy{}, errors.Wrap(err, fmt.Sprintf("error parsing content of policy %s", name))
	}
	return Policy{
		Name:    name,
		Content: content,
		Paths:   paths,
	}, nil
}
