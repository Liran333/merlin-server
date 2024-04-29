/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package messageadapter provides an adapter for working with message-related functionality.
package messageadapter

// Topics is a struct that represents the topics related to space deletion and update.
type Topics struct {
	SpaceCreated    string `json:"space_created" required:"true"`
	SpaceDeleted    string `json:"space_deleted" required:"true"`
	SpaceUpdated    string `json:"space_updated" required:"true"`
	SpaceEnvChanged string `json:"space_env_changed" required:"true"`
	SpaceDisable    string `json:"space_disable" required:"true"`
	SpaceForceEvent string `json:"space_force_event" required:"true"`
}
