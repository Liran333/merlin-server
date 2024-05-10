/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package messageadapter provides an adapter for working with message-related functionality.
package messageadapter

// Topics is a struct that represents the topics related to like deletion and creation.
type Topics struct {
	LikeCreate string `json:"like_create" required:"true"`
	LikeDelete string `json:"like_delete" required:"true"`
}
