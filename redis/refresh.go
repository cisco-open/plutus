package redis

import (
	"plutus/vault"
)

// RefreshUserToGroups refreshes the User to Group and the Group to User mappings in the database
func (c *Client) RefreshUserToGroups(data *vault.Data) error {
	namespace := data.Namespace
	c.deleteAllWithAnyPrefixes(namespace, PrefixUsrToVgr, PrefixVgrToUsr)

	groups := data.Groups
	aliases := data.Aliases

	for grpName, group := range groups {
		memberIDS := group.MemberEntityIds
		for _, mID := range memberIDS {
			username := aliases[mID].Name
			c.sAdd(namespace, PrefixUsrToVgr, username, grpName)
			c.sAdd(namespace, PrefixVgrToUsr, grpName, username)
		}
	}
	return nil
}

// RefreshGroupToPolicies refreshes the Group to Policy and the Policy to Group mappings in the database
func (c *Client) RefreshGroupToPolicies(data *vault.Data) error {
	namespace := data.Namespace
	c.deleteAllWithAnyPrefixes(namespace, PrefixVgrToPol, PrefixPolToVgr)

	groups := data.Groups
	for grpName, group := range groups {
		policies := group.Policies
		for _, policyName := range policies {
			c.sAdd(namespace, PrefixVgrToPol, grpName, policyName)
			c.sAdd(namespace, PrefixPolToVgr, policyName, grpName)
		}
	}

	return nil
}

// RefreshUserToEntity refreshes the Entity to Policy and the Policy To Entity mappinga in the database
func (c *Client) RefreshUserToEntity(data *vault.Data) error {
	namespace := data.Namespace
	c.deleteAllWithAnyPrefixes(namespace, PrefixUsrToEnt, PrefixEntToUsr)

	entities := data.Entities
	aliases := data.Aliases

	for entityID := range entities {
		alias := aliases[entityID].Name
		c.set(namespace, PrefixUsrToEnt, alias, entityID)
		c.set(namespace, PrefixEntToUsr, entityID, alias)
	}
	return nil
}

// RefreshEntityToPolicies refreshes the Entity to Policy mapping in the database
func (c *Client) RefreshEntityToPolicies(data *vault.Data) error {
	namespace := data.Namespace
	c.deleteAllWithAnyPrefixes(namespace, PrefixEntToPol, PrefixPolToEnt)

	entities := data.Entities

	for entityID, entity := range entities {
		policies := entity.Policies
		for _, policyName := range policies {
			c.sAdd(namespace, PrefixEntToPol, entityID, policyName)
			c.sAdd(namespace, PrefixPolToEnt, policyName, entityID)
		}
	}
	return nil
}

// RefreshUserToExternalGroups refreshes the User to ExternalGroupName mapping in the database
// Only those external groups are refreshed that are a part of the vault roles
func (c *Client) RefreshUserToExternalGroups(data *vault.Data) error {
	namespace := data.Namespace
	c.deleteAllWithAnyPrefixes(namespace, PrefixUsrToEgr, PrefixEgrToUsr)

	roles := data.Roles

	membersCache := make(map[string][]string) // group name to member names
	for _, role := range roles {
		externalGroup := role.ExternalGroupName

		var members []string
		if cachedMembers, ok := membersCache[externalGroup]; ok {
			members = cachedMembers
		} else {
			members, err := c.groupsReader.Members(externalGroup)
			if err != nil {
				return err
			}
			membersCache[externalGroup] = members
		}
		for _, member := range members {
			c.sAdd(namespace, PrefixUsrToEgr, member, externalGroup)
			c.sAdd(namespace, PrefixEgrToUsr, externalGroup, member)
		}
	}
	return nil
}

// RefreshExternalGroupToRoles refreshes the ExternalGroupName to Role mapping in the database
// Only those external groups are refrshed that are a part of the vault roles
func (c *Client) RefreshExternalGroupToRoles(data *vault.Data) error {
	namespace := data.Namespace
	c.deleteAllWithAnyPrefixes(namespace, PrefixEgrToRol, PrefixRolToEgr)

	roles := data.Roles

	for roleName, role := range roles {
		externalGroup := role.ExternalGroupName
		c.sAdd(namespace, PrefixEgrToRol, externalGroup, roleName)
		c.sAdd(namespace, PrefixRolToEgr, roleName, externalGroup)
	}
	return nil
}

// RefreshRoleToPolicies refreshes the Role to Policy mapping in the database
func (c *Client) RefreshRoleToPolicies(data *vault.Data) error {
	namespace := data.Namespace
	c.deleteAllWithAnyPrefixes(namespace, PrefixRolToPol, PrefixPolToRol)

	roles := data.Roles

	for roleName, role := range roles {
		policies := role.Policies
		for _, policy := range policies {
			c.sAdd(namespace, PrefixRolToPol, roleName, policy)
			c.sAdd(namespace, PrefixPolToRol, policy, roleName)
		}
	}
	return nil
}

// RefreshPathToPolicy refreshes the Path To Policy mapping in the database
func (c *Client) RefreshPathToPolicy(data *vault.Data) error {
	namespace := data.Namespace
	c.deleteAllWithAnyPrefixes(namespace, PrefixPatToPol, PrefixPolToPat)

	policies := data.Policies

	for policyName, policy := range policies {

		for _, path := range policy.Paths {

			encodedPath, err := path.Encode()
			if err != nil {
				return err
			}

			c.sAdd(namespace, PrefixPolToPat, policyName, encodedPath)
			// encodedPath is added to have the policy name and the capabilities
			c.sAdd(namespace, PrefixPatToPol, path.Path, encodedPath)

		}
	}

	return nil
}
