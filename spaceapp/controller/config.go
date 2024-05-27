/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package controller provides the controllers for handling HTTP requests and managing the application's business logic.
package controller

var config Config

// Init initializes the controller package with the provided configuration.
func Init(cfg *Config) {
	config = *cfg
}

// Config is a struct that holds the configuration for the controller package.
type Config struct {
	SSEToken string `json:"sse_token"`
}
