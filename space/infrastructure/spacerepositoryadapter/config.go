/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package spacerepositoryadapter provides an adapter for working with space repositories.
package spacerepositoryadapter

// Tables is a struct that represents a table with a space.
type Tables struct {
	Space string `json:"space" required:"true"`
}
