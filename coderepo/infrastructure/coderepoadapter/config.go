/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package coderepoadapter provides an adapter for interacting with a code repository service.
package coderepoadapter

// Config is a struct that represents the configuration for the code repository adapter.
type Config struct {
	ForceToBePrivate bool `json:"force_to_be_private"`
}
