package vault

import (
	"fmt"
	"plutus/parser"

	"github.com/pkg/errors"
)

// Policy is used to represent a vault policy
type Policy struct {
	Name    string
	Content string
	Paths   []parser.ParsedPath
}

// EncodedPaths returns the list of paths encoded as strings
// ready to be added to the redis instance
func (p *Policy) EncodedPaths() ([]string, error) {
	encodedPaths := make([]string, len(p.Paths))

	for i, path := range p.Paths {
		encodedPath, err := path.Encode()
		if err != nil {
			return nil, err
		}
		encodedPaths[i] = encodedPath
	}
	return encodedPaths, nil
}

// Policies returns a map of all policies on vault
// The mapping is from the policy name to the policy object
func (c *client) Policies() (map[string]Policy, error) {
	sys := c.client.Sys()

	policyNames, err := sys.ListPolicies()
	if err != nil {
		return nil, err
	}

	policies := make(map[string]Policy)

	for _, policyName := range policyNames {
		content, err := sys.GetPolicy(policyName)
		if err != nil {
			return nil, err
		}
		policies[policyName], err = newPolicy(policyName, content)
		if err != nil {
			c.logger.Warning(err)
		}
	}

	return policies, nil

}

func newPolicy(name, content string) (Policy, error) {
	hclParser := parser.NewHCLParser()
	paths, err := hclParser.ParseString(content, name)
	if err != nil {
		return Policy{}, errors.Wrap(err, fmt.Sprintf("error parsing content of policy %s", name))
	}
	return Policy{
		Name:    name,
		Content: content,
		Paths:   paths,
	}, nil
}
