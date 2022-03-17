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

package groupsreader

import (
	"context"
	"fmt"
	"os"

	b64 "encoding/base64"

	"plutus/config"
	"plutus/constants"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

// EnterpriseGithubGroupReader reads the groups from the repository specified in the config file
type EnterpriseGithubGroupReader struct {
	githubClient   *github.Client
	logger         *logrus.Logger
	groupsRepoPath string
}

var _ GroupsReader = &EnterpriseGithubGroupReader{}

// NewEnterpriseGithubGroupReader returns a NewEnterpriseGithubGroupReader initialised with a GitHub client
func NewEnterpriseGithubGroupReader(logger *logrus.Logger) (*EnterpriseGithubGroupReader, error) {
	githubToken := constants.ENV_GITHUB_TOKEN
	token, ok := os.LookupEnv(githubToken)
	if !ok {
		return nil, errors.Wrap(constants.EnvNotSetErr(githubToken), "error creating EnterpriseGithubGroupReader")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)

	config, err := config.GithubEnterpriseConfig()
	if err != nil {
		return nil, errors.Wrap(err, "error getting github enterprise reader base url")
	}
	baseURL := config.BaseURL()
	client, err := github.NewEnterpriseClient(baseURL, baseURL, tc)
	if err != nil {
		return nil, err
	}
	return &EnterpriseGithubGroupReader{
		githubClient:   client,
		logger:         logger,
		groupsRepoPath: config.GroupsRepoPath(),
	}, nil
}

// Members gets the alias of the members of the group grpName from the github repo
// If the group is not present on the repo, a list of empty members is returned
func (gr *EnterpriseGithubGroupReader) Members(grpName string) ([]string, error) {
	gr.logger.Infof("fetching group %s from Github", grpName)

	ctx := context.Background()

	requestURL := fmt.Sprintf("%s%s-groups.yaml", gr.groupsRepoPath, grpName)
	req, err := gr.githubClient.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error forming the request")
	}

	var f githubFile
	resp, err := gr.githubClient.Do(ctx, req, &f)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return []string{}, nil
		}
		return nil, errors.Wrap(err, "error making the request")
	}
	gr.logger.Info("remaining github requests for this hour: ", resp.Remaining)

	decodedContent, err := b64.StdEncoding.DecodeString(f.Content)

	var file yamlFile
	err = yaml.Unmarshal([]byte(decodedContent), &file)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling yaml")
	}
	return file.Members(), nil
}
