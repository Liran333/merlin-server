package app

import (
	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/domain"
	"github.com/openmerlin/merlin-server/models/domain/repository"
)

type CmdToCreateModel struct {
	coderepoapp.CmdToCreateRepo

	Desc     primitive.MSDDesc
	Fullname primitive.MSDFullname
}

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

	return
}

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

type ModelsDTO struct {
	Total  int                       `json:"total"`
	Models []repository.ModelSummary `json:"models"`
}

type CmdToListModels = repository.ListOption

type CmdToResetLabels = domain.ModelLabels
