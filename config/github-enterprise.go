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
	wrapErr "github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// GithubEnterpiseConfig is the config required by the github enterprise Groups Reader in the groups-reader package
type GithubEnterpiseConfig struct {
	GithubEnterpise struct {
		BaseURL        string `yaml:"baseURL"`
		GroupsRepoPath string `yaml:"groupsRepoPath"`
	} `yaml:"githubEnterpise"`
}

// BaseURL returns the baseURL of the github enterprise where the groups repository is stored
func (gec GithubEnterpiseConfig) BaseURL() string {
	return gec.GithubEnterpise.BaseURL
}

// GroupsRepoPath returns the repository path of the repository that has the groups yaml files
func (gec GithubEnterpiseConfig) GroupsRepoPath() string {
	return gec.GithubEnterpise.GroupsRepoPath
}

// GithubEnterpriseConfig the config for the github enterprise group reader
func GithubEnterpriseConfig() (GithubEnterpiseConfig, error) {
	fileBytes, err := readConfigFile()
	if err != nil {
		return GithubEnterpiseConfig{}, err
	}

	var githubEnterpriseConfig GithubEnterpiseConfig
	err = yaml.Unmarshal(fileBytes, &githubEnterpriseConfig)
	if err != nil {
		return GithubEnterpiseConfig{}, wrapErr.Wrapf(err, "error Unmarshaling config file %s", configFile)
	}

	if githubEnterpriseConfig.BaseURL() == "" {
		return GithubEnterpiseConfig{}, fieldNotSetError("githubEnterprise.baseURL")
	}

	if githubEnterpriseConfig.GroupsRepoPath() == "" {
		return GithubEnterpiseConfig{}, fieldNotSetError("githubEnterprise.groupsRepoPath")
	}
	return githubEnterpriseConfig, nil
}
