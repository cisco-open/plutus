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

package parser

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
)

// ParsedPath represents a path present in a vault policy
type ParsedPath struct {
	Path         string       `json:"path"`
	Capabilities Capabilities `json:"capabilities"`
	PolicyName   string       `json:"policy_name"`
}

// Encode returns the ParsedPath encoded as a string using base64 encoding
func (ps ParsedPath) Encode() (string, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	err := json.NewEncoder(encoder).Encode(ps)
	if err != nil {
		return "", err
	}
	encoder.Close()
	return buf.String(), nil

}

// DecodePath returns a ParsedPath decoded from the string using base64 encoding
func DecodePath(encodedStr string) (ParsedPath, error) {
	decoder := json.NewDecoder(base64.NewDecoder(base64.StdEncoding, strings.NewReader(encodedStr)))
	var ps ParsedPath
	err := decoder.Decode(&ps)

	if err != nil {
		return ParsedPath{}, err
	}

	return ps, nil
}

// DecodePaths returns a list ParsedPath decoded from the string using base64 encoding
func DecodePaths(encodedStrs []string) ([]ParsedPath, error) {

	parsedPaths := make([]ParsedPath, len(encodedStrs))

	for i, encodedStr := range encodedStrs {
		ps, err := DecodePath(encodedStr)
		if err != nil {
			return nil, err
		}
		parsedPaths[i] = ps
	}

	return parsedPaths, nil
}
