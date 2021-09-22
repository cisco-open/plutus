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
