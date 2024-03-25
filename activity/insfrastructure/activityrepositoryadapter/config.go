/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package modelrepositoryadapter provides an adapter for the model repository
// which handles database operations related to models.
package activityrepositoryadapter

// Tables is a struct that represents table names for different entities.
type Tables struct {
	Activity string `json:"activity" required:"true"`
}
