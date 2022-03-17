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

import "plutus/parser"

// PathResponse represents a path that has been fetched from redis
// It contains additional info about where the policy came from
type PathResponse struct {
	Path         string   `json:"path"`
	PolicyName   string   `json:"policy_name"`
	Capabilities []string `json:"capabilities"`
}

// NewPathResponse returns a new PathResponse
func NewPathResponse(path parser.ParsedPath) PathResponse {
	return PathResponse{
		Path:         path.Path,
		PolicyName:   path.PolicyName,
		Capabilities: path.Capabilities.Strings(),
	}
}

// PathResponseSet is a set of policy responses
type PathResponseSet map[string]PathResponse

// AsSlice returns all pathResponses as a slice
func (set PathResponseSet) AsSlice() []PathResponse {
	slice := make([]PathResponse, len(set))

	i := 0
	for _, pr := range set {
		slice[i] = pr
		i++
	}

	return slice
}
