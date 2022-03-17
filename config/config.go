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

package config

import (
	"io/ioutil"

	wrapErr "github.com/pkg/errors"

	"gopkg.in/yaml.v2"
)

const (
	configFile = "config/config.yaml"
)

// AppConfig is the config required by the application
type AppConfig struct {
	Namespaces []string `yaml:"namespaces"`
	UIAddress  string   `yaml:"uiAddress"`
}

func readConfigFile() ([]byte, error) {
	yFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, wrapErr.Wrapf(err, "error reading config file %s", configFile)
	}

	return yFile, nil
}

// appConfig returns the config file as a map
func appConfig() (AppConfig, error) {
	fileBytes, err := readConfigFile()
	if err != nil {
		return AppConfig{}, err
	}

	var appConfig AppConfig
	err = yaml.Unmarshal(fileBytes, &appConfig)
	if err != nil {
		return AppConfig{}, wrapErr.Wrapf(err, "Error Unmarshaling config file %s", configFile)
	}

	if appConfig.UIAddress == "" {
		return AppConfig{}, fieldNotSetError("uiAdress")
	}
	if appConfig.Namespaces == nil {
		return AppConfig{}, fieldNotSetError("namespaces")
	}

	return appConfig, nil
}

// GetNamespaces gets the uiAddress from the yaml config file
func GetNamespaces() ([]string, error) {
	config, err := appConfig()
	if err != nil {
		return nil, err
	}
	return config.Namespaces, nil
}

// GetUIAddress gets the uiAddress from the yaml config file
func GetUIAddress() (string, error) {
	config, err := appConfig()
	if err != nil {
		return "", err
	}
	return config.UIAddress, nil
}
