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

// Entity is used to represent a vault entity
type Entity struct {
	Name     string            `json:"name"`
	ID       string            `json:"id,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Policies []string          `json:"policies"`
	Disabled bool              `json:"disabled"`
}

func (c *client) Entities() (map[string]Entity, error) {
	entityIDS, err := c.client.Logical().List("identity/entity/id")
	if err != nil {
		return nil, err
	}

	keys := entityIDS.Data["keys"].([]interface{})

	entities := make(map[string]Entity)

	for _, key := range keys {
		id := key.(string)

		entity, err := c.getEntityByID(id)
		if err != nil {
			return nil, err
		}
		entities[entity.ID] = entity
	}

	return entities, nil
}

// GetEntityByID gets the entity by ID
func (c *client) getEntityByID(entityID string) (Entity, error) {
	r := c.client.NewRequest("GET", fmt.Sprintf("/v1/identity/entity/id/%s", entityID))

	resp, err := c.client.RawRequest(r)
	if err != nil {
		return Entity{}, errors.Wrap(err, "error doing http request")
	}

	type GetResponse struct {
		Data Entity
	}
	gResp := GetResponse{}
	err = resp.DecodeJSON(&gResp)
	if err != nil {
		return Entity{}, errors.Wrap(err, "error unmarshalling json body")
	}

	return gResp.Data, nil
}
