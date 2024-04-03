/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"
	"strings"
)

const (
	SortByMostLikes       = "most_likes"
	SortByAlphabetical    = "alphabetical"
	SortByMostDownloads   = "most_downloads"
	SortByRecentlyUpdated = "recently_updated"
	SortByRecentlyCreated = "recently_created"
	TrueCondition         = "1"
)

var (
	SortTypeMostLikes       = sortType(SortByMostLikes)
	SortTypeAlphabetical    = sortType(SortByAlphabetical)
	SortTypeMostDownloads   = sortType(SortByMostDownloads)
	SortTypeRecentlyUpdated = sortType(SortByRecentlyUpdated)
	SortTypeRecentlyCreated = sortType(SortByRecentlyCreated)
)

// SortType is an interface that defines a method to return the sort type as a string.
type SortType interface {
	SortType() string
}

// NewSortType creates a new SortType based on the input string.
func NewSortType(v string) (SortType, error) {
	switch strings.ToLower(v) {
	case SortByMostLikes:
		return SortTypeMostLikes, nil

	case SortByAlphabetical:
		return SortTypeAlphabetical, nil

	case SortByMostDownloads:
		return SortTypeMostDownloads, nil

	case SortByRecentlyUpdated:
		return SortTypeRecentlyUpdated, nil

	case SortByRecentlyCreated:
		return SortTypeRecentlyCreated, nil

	default:
		return nil, errors.New("unknown sort type")
	}
}

type sortType string

// SortType returns the sortType value as a string.
func (s sortType) SortType() string {
	return string(s)
}
