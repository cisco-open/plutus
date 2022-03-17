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

package constants

import "fmt"

// Environment variable lookup strings
const (
	ENV_REST_ADDR    = "REST_ADDR"
	ENV_REDIS_ADDR   = "REDIS_ADDR"
	ENV_GITHUB_TOKEN = "GITHUB_ACCESS_TOKEN"

	ENV_VAULT_ADDR      = "VAULT_ADDR"
	ENV_VAULT_TOKEN     = "VAULT_TOKEN"
	ENV_VAULT_NAMESPACE = "VAULT_NAMESPACE"
)

// EnvNotSetErr returns an error for when an ennvironment variable is s er
func EnvNotSetErr(envVar string) error {
	return fmt.Errorf("%s not set", envVar)
}
