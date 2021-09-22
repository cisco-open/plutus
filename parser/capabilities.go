package parser

import "strings"

// Capability is an enum for the capabilities defined in vault policies
type Capability int

// Capabilities possible on vault
const (
	Read   Capability = iota // Read   = 0
	Create                   // Create = 1
	Update                   // Update = 2
	Delete                   // Delete = 3
	List                     // List   = 4
	Sudo                     // Sudo   = 5
)

var capabilityStrs = []string{"read", "create", "update", "delete", "list", "sudo"}
var capabilitiesMap = map[string]Capability{
	"read":   Read,
	"create": Create,
	"update": Update,
	"delete": Delete,
	"list":   List,
	"sudo":   Sudo,
}

func (c Capability) String() string {
	return capabilityStrs[c]
}

// ParseCapability parses the string and returns the Capability
func ParseCapability(str string) (Capability, bool) {

	c, ok := capabilitiesMap[strings.ToLower(str)]
	return c, ok
}

// Capabilities is a list of capabilities
type Capabilities []Capability

// Strings returns Capabilities as a string slice
func (c Capabilities) Strings() []string {
	strs := make([]string, len(c))

	for i, capability := range c {
		strs[i] = capability.String()
	}

	return strs
}
