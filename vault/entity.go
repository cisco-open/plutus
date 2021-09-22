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
