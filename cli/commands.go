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

package cli

import (
	"github.com/urfave/cli/v2"
)

func (a *App) commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "refresh",
			Aliases: []string{"r"},
			Usage:   "Refreshes the data in the redis instance by updating it from vault",
			Action:  a.refreshAction,
		},
		{
			Name:    "start-rest-server",
			Aliases: []string{"s", "start"},
			Usage:   "Starts a REST API server on port 8000 unless another port is specified ",
			Action:  a.startAPIServer,
		},
		{
			Name:    "lookup",
			Aliases: []string{"l"},
			Usage:   "Allows you to lookup vault",
			Subcommands: []*cli.Command{
				{
					Name:      "policies",
					Usage:     "gets the policies for a given username",
					ArgsUsage: "[username]",
					Action:    a.lookupPoliciesAction,
				},
				{
					Name:      "users",
					Usage:     "gets the users for a given vault path",
					ArgsUsage: "[vault path or policy]",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:  "path",
							Value: false,
							Usage: "Lookup Users by vault path",
						},
						&cli.BoolFlag{
							Name:  "policy",
							Value: false,
							Usage: "Lookup Users by policies",
						},
					},
					Action: a.lookupUsersAction,
				},
			},
		},
	}
}
