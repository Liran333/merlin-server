/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package activityrepositoryadapter provides an adapter for the model repository
package activityrepositoryadapter

var config Config

// Tables is a struct that represents table names for different entities.
type Tables struct {
	Activity string `json:"activity" required:"true"`
}

// Config is a struct that holds the configuration for max number of activities a user can have.
type Config struct {
	MaxRecordPerPerson int64 `json:"max_record_per_person" required:"true"`
}

// InitUsage initializes the application with the provided configuration.
func InitUsage(cfg *Config) {
	config = *cfg
}
