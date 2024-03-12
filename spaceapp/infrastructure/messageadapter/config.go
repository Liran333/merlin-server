/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package messageadapter provides an adapter for working with message-related functionality.
package messageadapter

// Topics defines the topic names for message adapter operations.
type Topics struct {
	SpaceAppCreated   string `json:"space_app_created" required:"true"`
	SpaceAppRestarted string `json:"space_app_restarted" required:"true"`
}
