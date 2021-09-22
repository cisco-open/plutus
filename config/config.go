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
