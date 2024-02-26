/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a set of common primitive types and utilities.
package primitive

import "strings"

var (
	allSDK      map[string]bool
	allHardware map[string]bool
)

// Config represents the configuration structure for initialization.
type Config struct {
	SDK      []string `json:"sdk"      required:"true"`
	Hardware []string `json:"hardware" required:"true"`
}

// Init initializes the system with the provided configuration.
func Init(cfg *Config) {
	if cfg == nil {
		return
	}

	allHardware = map[string]bool{}
	if cfg.Hardware != nil {
		for _, v := range cfg.Hardware {
			allHardware[strings.ToLower(v)] = true
		}
	}

	allSDK = map[string]bool{}
	if cfg.SDK != nil {
		for _, sv := range cfg.SDK {
			allSDK[strings.ToLower(sv)] = true
		}
	}
}
