/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package messageadapter provides an adapter for working with message-related functionality.
package messageadapter

// Topics is a struct that represents the topics related to space deletion and update.
type Topics struct {
	SpaceDeleted string `json:"space_deleted" required:"true"`
	SpaceUpdated string `json:"space_updated" required:"true"`
}
