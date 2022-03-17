// Copyright 2022 Cisco Systems, Inc. and its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

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
