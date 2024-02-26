/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package primitive

var (
	maxBranchNameLength int
)

// Init initializes the configuration with the given Config struct.
func Init(cfg *Config) {
	maxBranchNameLength = cfg.MaxBranchNameLength
}

// Config represents the configuration for the application.
type Config struct {
	MaxBranchNameLength int `json:"max_branch_name_length" `
}

// SetDefault sets the default values for the Config struct if they are not set.
func (cfg *Config) SetDefault() {
	if cfg.MaxBranchNameLength <= 0 {
		cfg.MaxBranchNameLength = 100
	}
}
