/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

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
	MaxBranchNameLength int    `json:"max_branch_name_length"`
	BranchRegexp        string `json:"branch_regexp" required:"true"`

	branchRegexp *regexp.Regexp
}

// SetDefault sets the default values for the Config struct if they are not set.
func (cfg *Config) Validate() (err error) {
	cfg.branchRegexp, err = regexp.Compile(cfg.BranchRegexp)

	return
}
