package utils

import (
	"fmt"
	"strings"
)

// LookupString returns :
// 	NAMESPACE + "-" + PREFIX + KEY
// This is the the way the lookup string is formatted in the redis instance
func LookupString(namespace, prefix, key string) string {

	if strings.HasSuffix(namespace, "-") {
		panic(fmt.Sprintf("namespace %s should not end with '-'", namespace))
	}

	if !strings.HasSuffix(prefix, "-") {
		panic(fmt.Sprintf("prefix %s should end with '-'", prefix))
	}

	return fmt.Sprintf("%s-%s%s", namespace, prefix, key)
}

// NamespacePrefix rreturns the Namespace and the Prefix as a prefix
func NamespacePrefix(namespace, prefix string) string {
	return namespace + "-" + prefix
}
