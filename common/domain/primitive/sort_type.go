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
)

var (
	SortTypeMostLikes       = sortType(SortByMostLikes)
	SortTypeAlphabetical    = sortType(SortByAlphabetical)
	SortTypeMostDownloads   = sortType(SortByMostDownloads)
	SortTypeRecentlyUpdated = sortType(SortByRecentlyUpdated)
	SortTypeRecentlyCreated = sortType(SortByRecentlyCreated)
)

// SortType
type SortType interface {
	SortType() string
}

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

func (s sortType) SortType() string {
	return string(s)
}
