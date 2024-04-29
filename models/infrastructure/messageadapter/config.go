/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package messageadapter provides an adapter for working with message-related functionality.
package messageadapter

// Topics is a struct that represents the topics related to space deletion and update.
type Topics struct {
	ModelCreated string `json:"model_created" required:"true"`
	ModelUpdated string `json:"model_updated" required:"true"`
	ModelDeleted string `json:"model_deleted" required:"true"`
	ModelDisable string `json:"model_disable" required:"true"`
}
