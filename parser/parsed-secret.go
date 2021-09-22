package parser

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
)

// ParsedPath represents a path present in a vault policy
type ParsedPath struct {
	Path         string       `json:"path"`
	Capabilities Capabilities `json:"capabilities"`
	PolicyName   string       `json:"policy_name"`
}

// Encode returns the ParsedPath encoded as a string using base64 encoding
func (ps ParsedPath) Encode() (string, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	err := json.NewEncoder(encoder).Encode(ps)
	if err != nil {
		return "", err
	}
	encoder.Close()
	return buf.String(), nil

}

// DecodePath returns a ParsedPath decoded from the string using base64 encoding
func DecodePath(encodedStr string) (ParsedPath, error) {
	decoder := json.NewDecoder(base64.NewDecoder(base64.StdEncoding, strings.NewReader(encodedStr)))
	var ps ParsedPath
	err := decoder.Decode(&ps)

	if err != nil {
		return ParsedPath{}, err
	}

	return ps, nil
}

// DecodePaths returns a list ParsedPath decoded from the string using base64 encoding
func DecodePaths(encodedStrs []string) ([]ParsedPath, error) {

	parsedPaths := make([]ParsedPath, len(encodedStrs))

	for i, encodedStr := range encodedStrs {
		ps, err := DecodePath(encodedStr)
		if err != nil {
			return nil, err
		}
		parsedPaths[i] = ps
	}

	return parsedPaths, nil
}
