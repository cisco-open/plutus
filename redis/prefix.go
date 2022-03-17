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

package redis

import "plutus/utils"

// Prefixes for the keys in the redis instance
const (
	// USE CASE#1: USER -> POLICIES
	PrefixUsrToVgr = "usr2vgr-" // vault user to set(vault group names)
	PrefixVgrToPol = "vgr2pol-" // vault group name to set(policy names)

	PrefixUsrToEnt = "usr2ent-" // vault user to entityID
	PrefixEntToPol = "ent2pol-" // entityID to set(policy names)

	PrefixUsrToEgr = "usr2egr-" // vault user to set(external group names)
	PrefixEgrToRol = "egr2rol-" // external group name to set(vault role names)
	PrefixRolToPol = "rol2pol-" // role name to set(policy names)

	// USE CASE#2: PATH -> USERs
	PrefixPatToPol = "pat2pol-" // vault path to set(encoded policy name and capabilities)

	PrefixPolToVgr = "pol2vgr-" // vault policy name to set(vault group names)
	PrefixVgrToUsr = "vgr2usr-" // vault group name to set(vault users)

	PrefixPolToEnt = "pol2ent-" // policy names to set(entityIDs)
	PrefixEntToUsr = "ent2usr-" // vault entityID to user

	PrefixPolToRol = "pol2rol-" // policy name to set(role names)
	PrefixRolToEgr = "rol2egr-" // vault role name to set(external group names)
	PrefixEgrToUsr = "egr2usr-" // external group name to set(vault users)

	// MISC.
	PrefixPolToPat = "pol2sec-" // policy name to set(vault paths)
)

func (c *Client) deleteAllWithPrefix(prefix string) {
	iter := c.client.Scan(0, prefix+"*", 0).Iterator()
	for iter.Next() {
		c.client.Del(iter.Val())
	}
}

func (c *Client) deleteAllWithAnyPrefixes(namespace string, prefixes ...string) {
	for _, prefix := range prefixes {
		c.deleteAllWithPrefix(utils.NamespacePrefix(namespace, prefix))
	}
}
