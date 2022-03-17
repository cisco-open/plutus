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
