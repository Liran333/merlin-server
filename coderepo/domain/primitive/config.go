/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides primitive types and utility functions for working with basic concepts.
package primitive

import "regexp"

var (
	branchConfig Config
)

// Init initializes the configuration with the given Config struct.
func Init(cfg *Config) {
	branchConfig = *cfg
}

// Config represents the configuration for the application.
type Config struct {
	BranchRegexp        string `json:"branch_regexp"           required:"true"`
	BranchNameMinLength int    `json:"branch_name_min_length"  required:"true"`
	BranchNameMaxLength int    `json:"branch_name_max_length"  required:"true"`

	branchRegexp *regexp.Regexp
}

// Validate is a method that validates the Config instance.
func (cfg *Config) Validate() (err error) {
	cfg.branchRegexp, err = regexp.Compile(cfg.BranchRegexp)

	return
}
