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
