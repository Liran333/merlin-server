/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import "github.com/google/uuid"

// UUID is a type that represents a universally unique identifier.
type UUID = uuid.UUID

// NewUUID creates a new UUID from the given string.
func NewUUID(v string) (UUID, error) {
	return uuid.Parse(v)
}

// CreateUUID generates a new random UUID.
func CreateUUID() UUID {
	return uuid.New()
}
