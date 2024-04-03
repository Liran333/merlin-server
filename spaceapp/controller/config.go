/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package controller provides the controllers for handling HTTP requests and managing the application's business logic.
package controller

var config Config

func Init(cfg *Config) {
	config = *cfg
}

type Config struct {
	SSEToken string `json:"sse_token"`
}
