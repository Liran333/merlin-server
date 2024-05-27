/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package email provides functionality for sending emails.
package email

// Config represents the configuration for email.
type Config struct {
	AuthCode string `json:"auth_code" required:"true"`
	From     string `json:"from"      required:"true"`
	Host     string `json:"host"      required:"true"`
	Port     int    `json:"port"      required:"true"`
}
