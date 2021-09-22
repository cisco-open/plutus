package redis

import "plutus/utils"

// UserResponse represents the policy that has been fetched from redis
// It contains additional info about where the policy came from
type UserResponse struct {
	Username     string           `json:"username"`
	Policies     []PolicyResponse `json:"policies"`
	Capabilities []string         `json:"capabilities"`
}

// UserResponseSet is a set of User Responses
type UserResponseSet map[string]UserResponse

// AddSet adds the given otherSet to the first set
func (set UserResponseSet) AddSet(otherSet UserResponseSet) {
	for _, userResp := range otherSet {
		set.Add(userResp)
	}
}

// Add adds a UserResponse to the set
// If the policy already exists in the set,
// it is merged so as not to lose any info
func (set UserResponseSet) Add(cr UserResponse) {
	foundCr, ok := set[cr.Username]

	toAdd := cr
	if ok {
		foundPolicyResponses := NewPolicyResponseSet(foundCr.Policies)
		policyResponses := NewPolicyResponseSet(cr.Policies)

		mergedPolicyResponsesSet := policyResponses.Merge(foundPolicyResponses)

		foundCapabilitiesSet := utils.ToSet(foundCr.Capabilities)
		for _, capability := range cr.Capabilities {
			foundCapabilitiesSet[capability] = true
		}

		toAdd = UserResponse{
			Username:     cr.Username,
			Policies:     mergedPolicyResponsesSet.AsSlice(),
			Capabilities: utils.ToSlice(foundCapabilitiesSet),
		}
	}

	set[cr.Username] = toAdd
}

// AddPolicyResponse adds the giiven policyResponse to the UserResponce corresponding to the given username
func (set UserResponseSet) AddPolicyResponse(username string, pr PolicyResponse) {
	foundCr, ok := set[username]

	var toAdd UserResponse
	if ok {
		foundPoliciesSet := NewPolicyResponseSet(foundCr.Policies)
		foundPoliciesSet.Add(pr)

		foundCapabilitiesSet := utils.ToSet(foundCr.Capabilities)
		for _, capability := range pr.Capabilities {
			foundCapabilitiesSet[capability] = true
		}

		toAdd = UserResponse{
			Username:     username,
			Policies:     foundPoliciesSet.AsSlice(),
			Capabilities: utils.ToSlice(foundCapabilitiesSet),
		}
	} else {
		toAdd = UserResponse{
			Username:     username,
			Policies:     []PolicyResponse{pr},
			Capabilities: pr.Capabilities,
		}
	}

	set[username] = toAdd
}

// Users returns the list of users from the set
func (set UserResponseSet) Users() []string {
	users := make([]string, len(set))

	i := 0
	for user := range set {
		users[i] = user
		i++
	}

	return users
}

// AsSlice returns all policyResponses as a slice
func (set UserResponseSet) AsSlice() []UserResponse {
	slice := make([]UserResponse, len(set))

	i := 0
	for _, cr := range set {
		slice[i] = cr
		i++
	}

	return slice
}
