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

package main

import (
	"log"
	"plutus/cli"
	gr "plutus/groups-reader"

	"github.com/sirupsen/logrus"
)

// Logger used for logging
var logger = logrus.New()

func main() {
	logrus.SetLevel(logrus.ErrorLevel)

	groupsreader, err := gr.NewEnterpriseGithubGroupReader(logger)
	if err != nil {
		log.Fatal(err)
		logrus.Exit(1)
	}

	app, err := cli.NewApp(groupsreader, logger)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		log.Fatal(err) 
	}
}
