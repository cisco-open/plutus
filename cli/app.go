package cli

import (
	"os"
	gr "plutus/groups-reader"
	"plutus/redis"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// App is a CLI application
type App struct {
	cli          *cli.App
	redisClient  *redis.Client
	logger       *logrus.Logger
	groupsReader gr.GroupsReader
}

// NewApp returns a new CLI app
func NewApp(groupsReader gr.GroupsReader, logger *logrus.Logger) (*App, error) {

	redisClient, err := redis.NewClient(groupsReader, logger)
	if err != nil {
		return nil, err
	}

	app := &App{
		redisClient:  redisClient,
		logger:       logger,
		groupsReader: groupsReader,
	}

	cliApp := &cli.App{
		Name:     "plutus",
		Usage:    "A Vault Visualiser",
		Commands: app.commands(),
	}

	app.cli = cliApp
	return app, nil
}

// Run runs the CLI application
func (a *App) Run() error {
	return a.cli.Run(os.Args)
}
