/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package securestorage provides interfaces for defining secure manager for variable and secret.
package securestorage

// SpaceEnvSecret is an struct that represents an space store env secret.
type SpaceEnvSecret struct {
	Path  string
	Name  string
	Value string
}

// SpaceSecureManager is an interface that defines methods for sending space-related variable and secret.
type SpaceSecureManager interface {
	SaveSpaceEnvSecret(SpaceEnvSecret) error
	DeleteSpaceEnvSecret(string, string) error
	GetAllSpaceEnvSecret(SpaceEnvSecret) (string, error)
}
