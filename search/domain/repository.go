/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain models and configuration for a specific functionality.
package domain

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// Account is a type alias for primitive.Account.
type Account = primitive.Account

type SearchOption struct {
	SearchKey  string
	SearchType []string
	Account    Account
	Size       int
}

type SearchResult struct {
	SearchResultModel SearchResultModel `json:"model_result"`
	SearchResultSpace SearchResultSpace `json:"space_result"`
	SearchResultUser  SearchResultUser  `json:"user_result"`
	SearchResultOrg   SearchResultOrg   `json:"org_result"`
}

type SearchResultModel struct {
	ModelResult      []ModelResult `json:"result"`
	ModelResultCount int           `json:"count"`
}

type ModelResult struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
	Path  string `json:"path"`
}

type SearchResultSpace struct {
	SpaceResult      []SpaceResult `json:"result"`
	SpaceResultCount int           `json:"count"`
}

type SpaceResult struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
	Path  string `json:"path"`
}

type SearchResultOrg struct {
	OrgResult      []OrgResult `json:"result"`
	OrgResultCount int         `json:"count"`
}

type OrgResult struct {
	Id       string `json:"id"`
	Name     string `json:"account"`
	FullName string `json:"full_name"`
}

type SearchResultUser struct {
	UserResult      []UserResult `json:"result"`
	UserResultCount int          `json:"count"`
}

type UserResult struct {
	Account  string `json:"account"`
	FullName string `json:"full_name"`
	AvatarId string `json:"avatar_id"`
}
