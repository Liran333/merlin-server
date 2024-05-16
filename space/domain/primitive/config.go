/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a set of common primitive types and utilities.
package primitive

import (
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	envConfig  ENVConfig
	sdkObjects map[string]sets.Set[string]
	baseImages map[string]sets.Set[string]
	tasks      sets.Set[string]
)

// Config represents the configuration structure for initialization.
type Config struct {
	SDKObjects []SDKObject     `json:"sdk"`
	ENVConfig  ENVConfig       `json:"env"`
	Tasks      []string        `json:"tasks"      required:"true"`
	BaseImages []baseImageConf `json:"base_image"   required:"true"`
}

// ConfigItems returns a slice of interface{} containing pointers to the configuration items.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.ENVConfig,
	}
}

// ENVConfig represents the configuration for env.
type ENVConfig struct {
	MinValueLength int    `json:"env_value_min_length"      required:"true"`
	MaxValueLength int    `json:"env_value_max_length"      required:"true"`
	NameRegexp     string `json:"env_name_regexp"        required:"true"`

	nameRegexp *regexp.Regexp
}

// Validate check values for ENVConfig whether they are valid.
func (cfg *ENVConfig) Validate() (err error) {
	cfg.nameRegexp, err = regexp.Compile(cfg.NameRegexp)
	return
}

type SDKObject struct {
	SdkType  string   `json:"type"        required:"true"`
	Hardware []string `json:"hardware"    required:"true"`
}

type baseImageConf struct {
	HardwareType string   `json:"type"        required:"true"`
	BaseImage    []string `json:"base_image"    required:"true"`
}

// Init initializes the system with the provided configuration.
func Init(cfg *Config) {
	if cfg == nil {
		return
	}
	// init sdk
	sdkObjects = make(map[string]sets.Set[string])
	for _, sdkobj := range cfg.SDKObjects {
		sdkType := strings.ToLower(sdkobj.SdkType)
		for i, hardware := range sdkobj.Hardware {
			sdkobj.Hardware[i] = strings.ToLower(hardware)
		}
		sdkObjects[sdkType] = sets.New[string]()
		sdkObjects[sdkType].Insert(sdkobj.Hardware...)
	}

	envConfig = cfg.ENVConfig
	// init base image
	baseImages = make(map[string]sets.Set[string])
	for _, img := range cfg.BaseImages {
		hardwareType := strings.ToLower(img.HardwareType)
		for i, baseImage := range img.BaseImage {
			img.BaseImage[i] = strings.ToLower(baseImage)
		}
		baseImages[hardwareType] = sets.New[string]()
		baseImages[hardwareType].Insert(img.BaseImage...)
	}

	tasks = sets.New[string]()
	for _, task := range cfg.Tasks {
		tasks.Insert(task)
	}
}

// SetDefault sets default values for PasswordConfig if they are not provided.
func (cfg *Config) SetDefault() {
	if cfg.ENVConfig.MinValueLength <= 0 {
		cfg.ENVConfig.MinValueLength = 8
	}

	if cfg.ENVConfig.MaxValueLength <= 0 {
		cfg.ENVConfig.MaxValueLength = 20
	}
}
