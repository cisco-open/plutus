package config

import (
	"fmt"
)

func fieldNotSetError(fieldName string) error {
	return fmt.Errorf("%s field is not set in config file %s", fieldName, configFile)
}
