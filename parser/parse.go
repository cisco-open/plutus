package parser

import (
	errs "errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/pkg/errors"
)

// HCLParser is an adapter for hclparse.parser
type HCLParser struct {
	parser *hclparse.Parser
}

// NewHCLParser returns a new HCLParser
func NewHCLParser() *HCLParser {
	return &HCLParser{
		parser: hclparse.NewParser(),
	}
}

// ParseString parses the given HCL string and returns the list of vault paths in it
func (p *HCLParser) ParseString(hclStr, policyName string) ([]ParsedPath, error) {
	parser := p.parser

	bodyContent, diag := parser.ParseHCL([]byte(hclStr), policyName)
	if diag != nil {
		return nil, errors.Wrap(errs.New(diag.Error()), fmt.Sprintf("err parsing policy %s", policyName))
	}

	pathSchema := hcl.BlockHeaderSchema{
		Type:       "path",
		LabelNames: []string{""}, //TODO? What is its purpose (Why is it required)
	}

	hclSchema := hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{pathSchema},
	}

	// Parsing all the paths
	content, _ := bodyContent.Body.Content(&hclSchema)
	parsedPolicy := make([]ParsedPath, 0)

	for _, block := range content.Blocks {
		// Parsing the body of a path

		// Parsing the name
		pathName := block.Labels[0]
		attrs, _ := block.Body.JustAttributes()

		// Parsing the capabilities
		capabilitiesRaw, diag := attrs["capabilities"].Expr.Value(&hcl.EvalContext{})
		if diag != nil {
			return nil, errors.New(diag.Error())
		}

		capabilitiesSlice := capabilitiesRaw.AsValueSlice()
		capabilities := make([]Capability, len(capabilitiesSlice))
		for i, rawCapability := range capabilitiesSlice {
			ok := true
			capabilities[i], ok = ParseCapability(rawCapability.AsString())
			if !ok {
				return nil, fmt.Errorf("error in capabilities of path %s", pathName)
			}
		}

		path := ParsedPath{
			Path:         pathName,
			Capabilities: capabilities,
			PolicyName:   policyName,
		}

		parsedPolicy = append(parsedPolicy, path)
	}
	return parsedPolicy, nil
}
