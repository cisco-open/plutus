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
	"plutus/utils"
	"strings"
)

// PolicyResponse represents the policy that has been fetched from redis
// It contains additional info about where the policy came from
type PolicyResponse struct {
	Name         string   `json:"name"`
	OIDCRoles    []string `json:"roles,omitempty"`
	VaultGroups  []string `json:"groups,omitempty"`
	Entities     []string `json:"entities,omitempty"`
	Capabilities []string `json:"capabilities,omitempty"`
}

// NewPolicyResponse returns a new PolicyResponse
func NewPolicyResponse(name string) PolicyResponse {
	return PolicyResponse{
		Name:        name,
		OIDCRoles:   make([]string, 0),
		VaultGroups: make([]string, 0),
		Entities:    make([]string, 0),
	}
}

// Merge merges two PolicyResponse structs into one (if possible)
func (pr1 *PolicyResponse) Merge(pr2 PolicyResponse) (PolicyResponse, bool) {
	if pr1.Name != pr2.Name {
		return PolicyResponse{}, false
	}

	pr := PolicyResponse{
		Name:         pr1.Name,
		OIDCRoles:    utils.UniqueConcat(pr1.OIDCRoles, pr2.OIDCRoles...),
		VaultGroups:  utils.UniqueConcat(pr1.VaultGroups, pr2.VaultGroups...),
		Entities:     utils.UniqueConcat(pr1.Entities, pr2.Entities...),
		Capabilities: utils.UniqueConcat(pr1.Capabilities, pr2.Capabilities...),
	}

	return pr, true
}

// PolicyResponseSet is a set of policy responses
type PolicyResponseSet map[string]PolicyResponse

func (set PolicyResponseSet) String() string {
	sb := strings.Builder{}
	for policyName := range set {
		sb.WriteString(policyName + "\n")
	}
	return sb.String()
}

// NewPolicyResponseSet converts the given list of PolicyResponses into a set
func NewPolicyResponseSet(prs []PolicyResponse) PolicyResponseSet {
	set := make(PolicyResponseSet)

	for _, pr := range prs {
		set.Add(pr)
	}

	return set
}

// Merge merges 2 PolicyResponsesets
func (set PolicyResponseSet) Merge(otherSet PolicyResponseSet) PolicyResponseSet {
	res := set
	for _, policyResponse := range otherSet {
		res.Add(policyResponse)
	}
	return res
}

// Add adds a PolicyResponse to the set
// If the policy already exists in the set,
// it is merged so as not to lose any info
func (set PolicyResponseSet) Add(pr PolicyResponse) {
	policyResp, ok := set[pr.Name]

	// If policy exists in the set, merging it with the existing value
	var toAdd PolicyResponse
	if ok {
		mergedPr, _ := policyResp.Merge(pr)
		toAdd = mergedPr
	} else {
		toAdd = pr
	}

	set[pr.Name] = toAdd
}

// AsSlice returns all policyResponses as a slice
func (set PolicyResponseSet) AsSlice() []PolicyResponse {
	slice := make([]PolicyResponse, len(set))

	i := 0
	for _, pr := range set {
		slice[i] = pr
		i++
	}

	return slice
}
