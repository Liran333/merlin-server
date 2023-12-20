package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/domain"
)

type ModelSummary struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Desc          string `json:"desc"`
	Owner         string `json:"owner"`
	Fullname      string `json:"fullname"`
	TaskLabel     string `json:"task_label"`
	UpdatedAt     string `json:"updated_at"`
	LikeCount     int    `json:"like_count"`
	DownloadCount int    `json:"download_count"`
}

type ListOption struct {
	// can't define Name as domain.ResourceName
	// because the Name can be subpart of the real resource name
	Name string

	// list the models of Owner
	Owner primitive.Account

	// list by visibility
	Visibility primitive.Visibility

	SortType primitive.SortType

	// list models which have the labels
	Labels []string

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

type ModelRepositoryAdapter interface {
	Add(*domain.Model) error
	FindByName(primitive.Account, primitive.MSDName) (domain.Model, error)
	FindById(primitive.Identity) (domain.Model, error)
	Delete(primitive.Identity) error
	Save(*domain.Model) error
	List(*ListOption) ([]ModelSummary, int, error)
}
