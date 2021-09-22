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
