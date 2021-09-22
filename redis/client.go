package redis

import (
	"fmt"
	"log"
	"os"

	"plutus/constants"
	gr "plutus/groups-reader"
	"plutus/utils"
	"plutus/vault"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Client is an adapter for redis.Client
type Client struct {
	client       *redis.Client
	logger       *logrus.Logger
	groupsReader gr.GroupsReader
}

// sAdd calls the SADD command on the redis client
func (c *Client) sAdd(namespace, prefix, key, value string) {
	fullKey := utils.LookupString(namespace, prefix, key)
	c.logger.Infof("SADD %s %s", fullKey, value)
	c.client.SAdd(fullKey, value)
}

// get calls the GET command on the redis client
func (c *Client) get(namespace, prefix, key string) (string, error) {
	fullKey := utils.LookupString(namespace, prefix, key)
	c.logger.Infof("GET %s", fullKey)

	cmd := c.client.Get(fullKey)
	res, err := cmd.Result()
	c.logger.Infof(res)
	return res, err
}

// set calls the SET command on the redis client
func (c *Client) set(namespace, prefix, key, value string) {
	fullKey := utils.LookupString(namespace, prefix, key)
	c.logger.Infof("SET %s %s", fullKey, value)
	c.client.Set(fullKey, value, 0)
}

// smembers calls the SMEMBERS command on the redis client
// It tries 5 times with 5 second intervals and returns the error if it failes ultimately
func (c *Client) smembers(namespace, prefix, key string) (members []string, err error) {

	// function that actually calls smembers
	smembers := func() ([]string, error) {
		fullKey := utils.LookupString(namespace, prefix, key)
		c.logger.Infof("SMEMBERS %s", fullKey)
		cmd := c.client.SMembers(fullKey)
		res, err := cmd.Result()
		c.logger.Info(res)
		return res, err
	}

	// attempting to call smembers 5 times
	attempts := 5
	sleep := 5 * time.Second

	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Printf("attempt %d retrying after error: %v\n", i, err)
			time.Sleep(sleep)
			sleep *= 2
		}
		members, err = smembers()
		if err == nil {
			return members, nil
		}
	}

	return nil, fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

// NewClient returns a new Client object
func NewClient(groupsReader gr.GroupsReader, logger *logrus.Logger) (*Client, error) {
	redisAddr := constants.ENV_REDIS_ADDR
	address, ok := os.LookupEnv(redisAddr)
	if !ok {
		return nil, errors.Wrap(constants.EnvNotSetErr(redisAddr), "error creating redis client")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	if err := utils.DefaultRetry(client.Ping().Err); err != nil {
		return nil, errors.Wrap(err, "Error connecting to redis. PING failed")
	}

	return &Client{
		client:       client,
		logger:       logger,
		groupsReader: groupsReader,
	}, nil
}

// RefreshAll refreshes the different prefixed-defined "tables" in the redis instance concurrently
func (c *Client) RefreshAll(data *vault.Data) error {
	c.logger.Info("refreshing User to Groups")
	err := c.RefreshUserToGroups(data)
	if err != nil {
		return err
	}
	c.logger.Info("refreshing Group to Policies")
	err = c.RefreshGroupToPolicies(data)
	if err != nil {
		return err
	}

	c.logger.Info("refreshing User to Entity")
	err = c.RefreshUserToEntity(data)
	if err != nil {
		return err
	}
	c.logger.Info("refreshing Entity to Policies")
	err = c.RefreshEntityToPolicies(data)
	if err != nil {
		return err
	}

	c.logger.Info("refreshing User to External Groups")
	err = c.RefreshUserToExternalGroups(data)
	if err != nil {
		return err
	}
	c.logger.Info("refreshing External Group to Roles")
	err = c.RefreshExternalGroupToRoles(data)
	if err != nil {
		return err
	}
	c.logger.Info("refreshing Role to Policies")
	err = c.RefreshRoleToPolicies(data)
	if err != nil {
		return err
	}

	c.logger.Info("refreshing Path to (Encoded Policy and Path)")
	err = c.RefreshPathToPolicy(data)
	if err != nil {
		return err
	}

	return nil
}
