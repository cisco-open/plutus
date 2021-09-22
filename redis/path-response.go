package redis

import "plutus/parser"

// PathResponse represents a path that has been fetched from redis
// It contains additional info about where the policy came from
type PathResponse struct {
	Path         string   `json:"path"`
	PolicyName   string   `json:"policy_name"`
	Capabilities []string `json:"capabilities"`
}

// NewPathResponse returns a new PathResponse
func NewPathResponse(path parser.ParsedPath) PathResponse {
	return PathResponse{
		Path:         path.Path,
		PolicyName:   path.PolicyName,
		Capabilities: path.Capabilities.Strings(),
	}
}

// PathResponseSet is a set of policy responses
type PathResponseSet map[string]PathResponse

// AsSlice returns all pathResponses as a slice
func (set PathResponseSet) AsSlice() []PathResponse {
	slice := make([]PathResponse, len(set))

	i := 0
	for _, pr := range set {
		slice[i] = pr
		i++
	}

	return slice
}
