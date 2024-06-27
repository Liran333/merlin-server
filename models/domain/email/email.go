/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package email provides functionality for sending emails.
package email

type Email interface {
	Send(modelName, content, user, url string) error
	GetRootUrl() string
}
