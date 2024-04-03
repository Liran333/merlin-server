/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package branchrepositoryadapter provides an adapter for the branch repository using GORM.
package branchrepositoryadapter

// Tables is a struct that represents table names for different entities.
type Tables struct {
	Branch string `json:"branch" required:"true"`
}
