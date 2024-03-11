/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a set of common primitive types and utilities.
package primitive

import (
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	sdkObjects map[string]sets.Set[string]
)

// Config represents the configuration structure for initialization.
type Config struct {
	SDKObjects []SDKObject `json:"sdk"`
}

type SDKObject struct {
	SdkType  string   `json:"type"        required:"true"`
	Hardware []string `json:"hardware"    required:"true"`
}

// Init initializes the system with the provided configuration.
func Init(cfg *Config) {
	if cfg == nil {
		return
	}

	sdkObjects = make(map[string]sets.Set[string])
	for _, sdkobj := range cfg.SDKObjects {
		sdkType := strings.ToLower(sdkobj.SdkType)
		for i, hardware := range sdkobj.Hardware {
			sdkobj.Hardware[i] = strings.ToLower(hardware)
		}
		sdkObjects[sdkType] = sets.New[string]()
		sdkObjects[sdkType].Insert(sdkobj.Hardware...)
	}
}
