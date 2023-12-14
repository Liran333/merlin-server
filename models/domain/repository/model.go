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

	// LastId is id of the last element on the previous page
	LastId       string
	PageNum      int // will return the total num when PageNum == 1
	CountPerPage int
}

type ModelRepositoryAdapter interface {
	Add(*domain.Model) error
	FindByName(primitive.Account, primitive.MSDName) (domain.Model, error)
	//FindById(string)
	Delete(string) error
	Save(*domain.Model) error
	List(*ListOption) ([]ModelSummary, int, error)
}
