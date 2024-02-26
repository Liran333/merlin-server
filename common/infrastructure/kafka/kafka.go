/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package kafka provides functionality for interacting with Kafka.
package kafka

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/kafka-lib/mq"
)

const (
	deaultVersion = "2.1.0"
)

// Exit is an exported variable that provides the exit function for the Kafka package.
var Exit = kfklib.Exit

// Config represents the configuration for Kafka.
type Config struct {
	kfklib.Config
}

// SetDefault sets the default values for the Config.
func (cfg *Config) SetDefault() {
	if cfg.Version == "" {
		cfg.Version = deaultVersion
	}
}

// Init initializes the Kafka agent with the specified configuration, logger, and removeCfg flag.
func Init(cfg *Config, log mq.Logger, removeCfg bool) error {
	return kfklib.Init(&cfg.Config, log, nil, "", removeCfg)
}
