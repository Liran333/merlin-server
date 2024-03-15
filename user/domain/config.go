/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain models and configuration for a specific functionality.
package domain

type tables struct {
	User  string `json:"user"  required:"true"`
	Token string `json:"token" required:"true"`
}

// Config is a struct that holds the configuration for the program.
type Config struct {
	Tables tables `json:"tables" required:"true"`
	Key    []byte `json:"key"    required:"true"`
}
