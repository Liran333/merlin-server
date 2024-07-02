/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package modelrepositoryadapter provides an adapter for the model repository
// which handles database operations related to models.
package modelrepositoryadapter

// Tables is a struct that represents table names for different entities.
type Tables struct {
	Model       string `json:"model" required:"true"`
	ModelDeploy string `json:"model_deploy" required:"true"`
}
