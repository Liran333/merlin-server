/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repositoryadapter provides an adapter implementation for working with the repository of space applications.
package repositoryadapter

// Tables is a struct that represents table names for different entities.
type Tables struct {
	SpaceApp string `json:"space_app" required:"true"`
}
