/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package loginrepositoryadapter provides a data structure for the "Tables" type used in the login repository adapter.
package loginrepositoryadapter

// Tables represents the tables used in the login repository adapter.
type Tables struct {
	Login string `json:"login" required:"true"`
}
