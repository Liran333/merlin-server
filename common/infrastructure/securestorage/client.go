/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package securestorage provides interfaces for defining secure manager for variable and secret.
package securestorage

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/userpass"
	"github.com/sirupsen/logrus"
)

var (
	cli *api.Client
)

// Init initializes the vault client with the given configuration.
func Init(config *Config) error {
	// init vault client
	defaultConfig := api.DefaultConfig()

	defaultConfig.Address = config.Address

	// init http.Transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // #nosec G402
		},
	}
	defaultConfig.HttpClient.Transport = tr

	client, err := api.NewClient(defaultConfig)
	if err != nil {
		logrus.Errorf("unable to initialize Vault client: %v", err)
		return err
	}
	userpassAuth, err := userpass.NewUserpassAuth(config.UserName, &userpass.Password{FromString: config.PassWord})
	if err != nil {
		logrus.Errorf("initialize vault userpass auth failed: %v", err)
		return err
	}
	loginRespFromFile, err := client.Auth().Login(context.Background(), userpassAuth)
	if err != nil {
		logrus.Errorf("unable to initialize userpass auth method: %v", err)
		return err
	}
	if loginRespFromFile.Auth == nil || loginRespFromFile.Auth.ClientToken == "" {
		logrus.Errorf("unable to initialize userpass auth method: %v", err)
		return err
	}
	cli = client
	return nil
}

// GetClient returns the vault client.
func GetClient() *api.Client {
	return cli
}
