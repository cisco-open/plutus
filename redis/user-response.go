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

import "plutus/utils"

// UserResponse represents the policy that has been fetched from redis
// It contains additional info about where the policy came from
type UserResponse struct {
	Username     string           `json:"username"`
	Policies     []PolicyResponse `json:"policies"`
	Capabilities []string         `json:"capabilities"`
}

// UserResponseSet is a set of User Responses
type UserResponseSet map[string]UserResponse

// AddSet adds the given otherSet to the first set
func (set UserResponseSet) AddSet(otherSet UserResponseSet) {
	for _, userResp := range otherSet {
		set.Add(userResp)
	}
}

// Add adds a UserResponse to the set
// If the policy already exists in the set,
// it is merged so as not to lose any info
func (set UserResponseSet) Add(cr UserResponse) {
	foundCr, ok := set[cr.Username]

	toAdd := cr
	if ok {
		foundPolicyResponses := NewPolicyResponseSet(foundCr.Policies)
		policyResponses := NewPolicyResponseSet(cr.Policies)

		mergedPolicyResponsesSet := policyResponses.Merge(foundPolicyResponses)

		foundCapabilitiesSet := utils.ToSet(foundCr.Capabilities)
		for _, capability := range cr.Capabilities {
			foundCapabilitiesSet[capability] = true
		}

		toAdd = UserResponse{
			Username:     cr.Username,
			Policies:     mergedPolicyResponsesSet.AsSlice(),
			Capabilities: utils.ToSlice(foundCapabilitiesSet),
		}
	}

	set[cr.Username] = toAdd
}

// AddPolicyResponse adds the giiven policyResponse to the UserResponce corresponding to the given username
func (set UserResponseSet) AddPolicyResponse(username string, pr PolicyResponse) {
	foundCr, ok := set[username]

	var toAdd UserResponse
	if ok {
		foundPoliciesSet := NewPolicyResponseSet(foundCr.Policies)
		foundPoliciesSet.Add(pr)

		foundCapabilitiesSet := utils.ToSet(foundCr.Capabilities)
		for _, capability := range pr.Capabilities {
			foundCapabilitiesSet[capability] = true
		}

		toAdd = UserResponse{
			Username:     username,
			Policies:     foundPoliciesSet.AsSlice(),
			Capabilities: utils.ToSlice(foundCapabilitiesSet),
		}
	} else {
		toAdd = UserResponse{
			Username:     username,
			Policies:     []PolicyResponse{pr},
			Capabilities: pr.Capabilities,
		}
	}

	set[username] = toAdd
}

// Users returns the list of users from the set
func (set UserResponseSet) Users() []string {
	users := make([]string, len(set))

	i := 0
	for user := range set {
		users[i] = user
		i++
	}

	return users
}

// AsSlice returns all policyResponses as a slice
func (set UserResponseSet) AsSlice() []UserResponse {
	slice := make([]UserResponse, len(set))

	i := 0
	for _, cr := range set {
		slice[i] = cr
		i++
	}

	return slice
}
