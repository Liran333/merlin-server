/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package messageadapter provides an adapter for working with message-related functionality.
package messageadapter

// Topics is a struct that represents the topics related to dataset deletion and update.
type Topics struct {
	DatasetCreated string `json:"dataset_created" required:"true"`
	DatasetUpdated string `json:"dataset_updated" required:"true"`
	DatasetDeleted string `json:"dataset_deleted" required:"true"`
}
