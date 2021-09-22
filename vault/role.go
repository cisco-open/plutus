package vault

import (
	"regexp"
)

// Role is used to represent an OIDC role on Vault
type Role struct {
	Name           string
	ExternalGroupName string
	Policies       []string
}

const (
	strBeforeGrpName = "CN="
	strAfterGrpName  = ",OU"
)

func (c *client) Roles() (map[string]Role, error) {
	roleNames, err := c.client.Logical().List("auth/oidc/role")
	if err != nil {
		return nil, err
	}

	keys := roleNames.Data["keys"].([]interface{})
	roles := make(map[string]Role)

	for _, key := range keys {
		name := key.(string)
		role, err := c.getRoleByName(name)
		if err != nil {
			return nil, err
		}
		roles[name] = role
	}

	return roles, nil
}

func (c *client) getRoleByName(name string) (Role, error) {
	role, err := c.client.Logical().Read("auth/oidc/role/" + name)
	if err != nil {
		return Role{}, err
	}

	policyNamesFetched := role.Data["policies"]

	var policyNames []string
	if policyNamesFetched != nil {
		policyNamesRaw := policyNamesFetched.([]interface{})
		policyNames = make([]string, len(policyNamesRaw))
		for i, policyName := range policyNamesRaw {
			policyNames[i] = policyName.(string)
		}
	} else {
		c.logger.Warningf("role %s does not have any policies attached to it", name)
	}

	var externalGroup string
	if boundClaims := role.Data["bound_claims"]; boundClaims != nil {
		if memberOf := boundClaims.(map[string]interface{})["memberof"]; memberOf != nil {
			externalGroup = memberOf.(string)
			re := regexp.MustCompile(strBeforeGrpName + "(.)*" + strAfterGrpName)
			matched := re.FindString(externalGroup)
			externalGroup = matched[len(strBeforeGrpName) : len(matched)-len(strAfterGrpName)]
		}
	}

	return Role{
		Name:           name,
		ExternalGroupName: externalGroup,
		Policies:       policyNames,
	}, nil
}
