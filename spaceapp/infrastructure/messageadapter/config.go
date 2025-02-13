/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package messageadapter provides an adapter for working with message-related functionality.
package messageadapter

// Topics defines the topic names for message adapter operations.
type Topics struct {
	SpaceAppCreated   	string `json:"space_app_created" required:"true"`
	SpaceAppRestarted 	string `json:"space_app_restarted" required:"true"`
	SpaceAppPaused    	string `json:"space_app_paused" required:"true"`
	SpaceAppResumed   	string `json:"space_app_resumed" required:"true"`
	SpaceAppHeartbeat   string `json:"space_app_heartbeat" required:"true"`
	SpaceAppSleep   	string `json:"space_app_sleep" required:"true"`
	SpaceAppWakeup   	string `json:"space_app_wakeup" required:"true"`
	SpaceForceEvent   	string `json:"space_force_event" required:"true"`
}
