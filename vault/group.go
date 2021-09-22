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
