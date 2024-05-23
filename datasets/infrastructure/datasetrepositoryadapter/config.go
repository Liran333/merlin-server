/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package datasetrepositoryadapter provides an adapter for the datasets repository
// which handles database operations related to datasets.
package datasetrepositoryadapter

// Tables is a struct that represents table names for different entities.
type Tables struct {
	Datasets string `json:"datasets" required:"true"`
}
