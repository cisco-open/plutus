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
	"net/http"

	"github.com/pkg/errors"
)

// Group is used to represent a vault group
type Group struct {
	ID              string            `json:"id,omitempty"`
	Name            string            `json:"name"`
	Type            string            `json:"type"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	Policies        []string          `json:"policies,omitempty"`
	MemberGroupIds  []string          `json:"member_group_ids,omitempty"`
	MemberEntityIds []string          `json:"member_entity_ids,omitempty"`
}

// Groups returns the list of all groups on vault
func (c *client) Groups() (map[string]Group, error) {
	groupIDS, err := c.client.Logical().List("identity/group/id")
	if err != nil {
		return nil, err
	}

	groups := make(map[string]Group)

	if groupIDS == nil {
		return groups, nil
	}

	keysRaw := groupIDS.Data["keys"]

	if keysRaw == nil {
		return groups, err
	}
	keys := keysRaw.([]interface{})

	for _, key := range keys {
		id := key.(string)

		group, err := c.getGroupByID(id)
		if err != nil {
			return nil, err
		}
		groups[group.Name] = group
	}

	return groups, nil
}

func (c *client) getGroupByID(id string) (Group, error) {
	r := c.client.NewRequest("GET", fmt.Sprintf("/v1/identity/group/id/%s", id))
	resp, err := c.client.RawRequest(r)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return Group{}, errors.Wrap(err, "group not found")
		}
		return Group{}, errors.Wrap(err, "error doing http request")
	}

	type GetResponse struct {
		Data Group
	}

	gResp := GetResponse{}
	err = resp.DecodeJSON(&gResp)
	if err != nil {
		return Group{}, errors.Wrap(err, "error unmarshalling json body")
	}

	return gResp.Data, nil
}
