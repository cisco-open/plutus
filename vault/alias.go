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

	"github.com/pkg/errors"
)

// Alias is used to represent the vault alias of a vault entity
type Alias struct {
	ID            string `json:"id,omitempty"`
	EntityID      string `json:"canonical_id"`
	MountAccessor string `json:"mount_accessor"`
	MountPath     string `json:"mount_path"`
	MountType     string `json:"mount_type"`
	Name          string `json:"name"`
}

func (c *client) Aliases() (map[string]Alias, error) {
	r := c.client.NewRequest("LIST", "/v1/identity/entity-alias/id")
	resp, err := c.client.RawRequest(r)
	if err != nil {
		return map[string]Alias{}, errors.Wrap(err, "error doing http request")
	}

	type ResponseData struct {
		KeyInfo map[string]Alias `json:"key_info"`
	}
	type ListResponse struct {
		Data ResponseData
	}
	lResp := ListResponse{}
	err = resp.DecodeJSON(&lResp)
	if err != nil {
		return map[string]Alias{}, errors.Wrap(err, "error unmarshalling json body")
	}

	aliasMap := make(map[string]Alias)
	for k, v := range lResp.Data.KeyInfo {
		v.ID = k
		aliasMap[v.EntityID] = v
	}

	return aliasMap, nil

}

func (c *client) GetAliasByEntityID(entityID string) (Alias, error) {

	r := c.client.NewRequest("GET", fmt.Sprintf("/v1/identity/entity-alias/id/%s", entityID))

	resp, err := c.client.RawRequest(r)
	if err != nil {
		return Alias{}, errors.Wrap(err, "error making http request")
	}

	type GetResponse struct {
		Data Alias
	}

	gResp := GetResponse{}
	err = resp.DecodeJSON(&gResp)
	if err != nil {
		return Alias{}, errors.Wrap(err, "error unmarshalling json body")
	}

	return gResp.Data, nil
}
