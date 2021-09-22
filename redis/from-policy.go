package redis

// QueryUsersFromPolicyName gets the Users that have the provided policy
func (c *Client) QueryUsersFromPolicyName(namespace, policyName string) ([]string, error) {
	usersSet := make(map[string]bool, 0)

	groups, err := c.smembers(namespace, PrefixPolToVgr, policyName)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		grpUsers, err := c.smembers(namespace, PrefixVgrToUsr, group)
		if err != nil {
			return nil, err
		}

		for _, grpUser := range grpUsers {
			usersSet[grpUser] = true
		}
	}

	entities, err := c.smembers(namespace, PrefixPolToEnt, policyName)
	if err != nil {
		return nil, err
	}

	for _, entityID := range entities {
		username, err := c.get(namespace, PrefixEntToUsr, entityID)
		if err != nil {
			return nil, err
		}
		usersSet[username] = true

	}

	roles, err := c.smembers(namespace, PrefixPolToRol, policyName)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		externalGroups, err := c.smembers(namespace, PrefixRolToEgr, role)
		if err != nil {
			return nil, err
		}

		for _, externalGroup := range externalGroups {
			externalUsers, err := c.smembers(namespace, PrefixEgrToUsr, externalGroup)
			if err != nil {
				return nil, err
			}

			for _, externalUser := range externalUsers {
				usersSet[externalUser] = true
			}
		}
	}

	i := 0
	usersSlice := make([]string, len(usersSet))
	for username := range usersSet {
		usersSlice[i] = username
		i++
	}

	return usersSlice, nil
}
