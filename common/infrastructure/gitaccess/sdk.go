/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package gitaccess the client of gitaccess.
package gitaccess

import (
	"github.com/openmerlin/git-access-sdk/httpclient"
)

func Init(cfg *Config) {
	httpclient.Init(cfg)
}

// Config is for http client config
type Config = httpclient.Config
