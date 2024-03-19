/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/domain"
	"github.com/openmerlin/merlin-server/models/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

// CmdToCreateModel is a struct that represents a command to create a model.
type CmdToCreateModel struct {
	coderepoapp.CmdToCreateRepo

	Desc     primitive.MSDDesc
	Fullname primitive.MSDFullname
}

// CmdToUpdateModel is a struct that represents a command to update a model.
type CmdToUpdateModel struct {
	coderepoapp.CmdToUpdateRepo

	Desc     primitive.MSDDesc
	Fullname primitive.MSDFullname
}

func (cmd *CmdToUpdateModel) toModel(model *domain.Model) (b bool) {
	if v := cmd.Desc; v != nil && v != model.Desc {
		model.Desc = v
		b = true
	}

	if v := cmd.Fullname; v != nil && v != model.Fullname {
		model.Fullname = v
		b = true
	}

	if b {
		model.UpdatedAt = utils.Now()
	}

	return
}

// ModelDTO is a struct that represents a data transfer object for a model.
type ModelDTO struct {
	Id            string         `json:"id"`
	Name          string         `json:"name"`
	Desc          string         `json:"desc"`
	Owner         string         `json:"owner"`
	Labels        ModelLabelsDTO `json:"labels"`
	Fullname      string         `json:"fullname"`
	CreatedAt     int64          `json:"created_at"`
	UpdatedAt     int64          `json:"updated_at"`
	LikeCount     int            `json:"like_count"`
	Visibility    string         `json:"visibility"`
	DownloadCount int            `json:"download_count"`
}

// ModelLabelsDTO is a struct that represents a data transfer object for model labels.
type ModelLabelsDTO struct {
	Task       string   `json:"task"`
	Others     []string `json:"others"`
	License    string   `json:"license"`
	Frameworks []string `json:"frameworks"`
}

func toModelLabelsDTO(model *domain.Model) ModelLabelsDTO {
	labels := &model.Labels

	return ModelLabelsDTO{
		Task:       labels.Task,
		Others:     labels.Others.UnsortedList(),
		License:    model.License.License(),
		Frameworks: labels.Frameworks.UnsortedList(),
	}
}

func toModelDTO(model *domain.Model) ModelDTO {
	dto := ModelDTO{
		Id:            model.Id.Identity(),
		Name:          model.Name.MSDName(),
		Owner:         model.Owner.Account(),
		Labels:        toModelLabelsDTO(model),
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
		LikeCount:     model.LikeCount,
		Visibility:    model.Visibility.Visibility(),
		DownloadCount: model.DownloadCount,
	}

	if model.Desc != nil {
		dto.Desc = model.Desc.MSDDesc()
	}

	if model.Fullname != nil {
		dto.Fullname = model.Fullname.MSDFullname()
	}

	return dto
}

// ModelsDTO is a struct that represents a data transfer object for a list of models.
type ModelsDTO struct {
	Total  int                       `json:"total"`
	Models []repository.ModelSummary `json:"models"`
}

// CmdToListModels is a type alias for repository.ListOption, representing a command to list models.
type CmdToListModels = repository.ListOption

// CmdToResetLabels is a type alias for domain.ModelLabels, representing a command to reset model labels.
type CmdToResetLabels = domain.ModelLabels
