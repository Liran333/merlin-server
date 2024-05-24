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

// SearchOption is a type alias for SearchOption.
type SearchOption struct {
	SearchKey  string
	SearchType []string
	Account    Account
	Size       int
}

// SearchResult is a type alias for SearchResult.
type SearchResult struct {
	SearchResultModel   SearchResultModel   `json:"model_result"`
	SearchResultDataset SearchResultDataset `json:"dataset_result"`
	SearchResultSpace   SearchResultSpace   `json:"space_result"`
	SearchResultUser    SearchResultUser    `json:"user_result"`
	SearchResultOrg     SearchResultOrg     `json:"org_result"`
}

// SearchResultModel is a type alias for SearchResultModel.
type SearchResultModel struct {
	ModelResult      []ModelResult `json:"result"`
	ModelResultCount int           `json:"count"`
}

// ModelResult is a type alias for ModelResult.
type ModelResult struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
	Path  string `json:"path"`
}

// SearchResultDataset represents the search result for datasets.
type SearchResultDataset struct {
	DatasetResult      []DatasetResult `json:"result"`
	DatasetResultCount int             `json:"count"`
}

// DatasetResult represents a single dataset in the search result.
type DatasetResult struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
	Path  string `json:"path"`
}

// SearchResultSpace represents the search result for spaces.
type SearchResultSpace struct {
	SpaceResult      []SpaceResult `json:"result"`
	SpaceResultCount int           `json:"count"`
}

// SpaceResult represents a single space in the search result.
type SpaceResult struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
	Path  string `json:"path"`
}

// SearchResultOrg represents the search result for orgs.
type SearchResultOrg struct {
	OrgResult      []OrgResult `json:"result"`
	OrgResultCount int         `json:"count"`
}

// OrgResult represents a single org in the search result.
type OrgResult struct {
	Id       string `json:"id"`
	Name     string `json:"account"`
	FullName string `json:"full_name"`
}

// SearchResultUser represents the search result for users.
type SearchResultUser struct {
	UserResult      []UserResult `json:"result"`
	UserResultCount int          `json:"count"`
}

// UserResult represents a single user in the search result.
type UserResult struct {
	Account  string `json:"account"`
	FullName string `json:"full_name"`
	AvatarId string `json:"avatar_id"`
}
