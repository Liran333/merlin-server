package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space/domain"
)

type SpaceSummary struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Desc          string `json:"desc"`
	Task          string `json:"task"`
	Owner         string `json:"owner"`
	Fullname      string `json:"fullname"`
	UpdatedAt     int64  `json:"updated_at"`
	LikeCount     int    `json:"like_count"`
	DownloadCount int    `json:"download_count"`
}

type ListOption struct {
	// can't define Name as domain.ResourceName
	// because the Name can be subpart of the real resource name
	Name string

	// list the space of Owner
	Owner primitive.Account

	// list by visibility
	Visibility primitive.Visibility

	// list space which have one of licenses
	License primitive.License

	// list space which have at least one label for each kind of lables.
	Labels domain.SpaceLabels

	// sort
	SortType primitive.SortType

	// whether to calculate the total
	Count        bool
	PageNum      int
	CountPerPage int
}

func (opt *ListOption) Pagination() (bool, int) {
	if opt.PageNum > 0 && opt.CountPerPage > 0 {
		return true, (opt.PageNum - 1) * opt.CountPerPage
	}

	return false, 0
}

type SpaceRepositoryAdapter interface {
	Add(*domain.Space) error
	FindByName(*domain.SpaceIndex) (domain.Space, error)
	FindById(primitive.Identity) (domain.Space, error)
	Delete(primitive.Identity) error
	Save(*domain.Space) error
	List(*ListOption) ([]SpaceSummary, int, error)
}

type SpaceLabelsRepoAdapter interface {
	Save(*domain.SpaceIndex, *domain.SpaceLabels) error
}
