/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain space app and configuration for the app service.
package domain

const (
	overRestartTimePeriod = 60 * 60 * 2
	overResumeTimePeriod  = 60 * 60 * 2
)

// Init initializes the configuration with the given Config struct.
var config Config

func Init(cfg *Config) {
	config = *cfg
}

// Config is a struct that holds the configuration for over restart time,
type Config struct {
	RestartOverTime int64 `json:"restart_over_time"`
	ResumeOverTime  int64 `json:"resume_over_time"`
}

// SetDefault sets the default values for the Config struct.
func (cfg *Config) SetDefault() {
	if cfg.RestartOverTime <= 0 {
		cfg.RestartOverTime = overRestartTimePeriod
	}
	if cfg.ResumeOverTime <= 0 {
		cfg.ResumeOverTime = overResumeTimePeriod
	}
}
