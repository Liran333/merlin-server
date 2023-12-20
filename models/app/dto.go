package app

import (
	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/domain"
	"github.com/openmerlin/merlin-server/models/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

type CmdToCreateModel struct {
	coderepoapp.CmdToCreateRepo

	Desc     primitive.MSDDesc
	Fullname primitive.MSDFullname
}

type ModelIndex struct {
	Owner primitive.Account
	Name  primitive.MSDName
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
	Id            string   `json:"id"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Owner         string   `json:"owner"`
	Labels        []string `json:"labels"`
	License       string   `json:"license"`
	Fullname      string   `json:"fullname"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
	LikeCount     int      `json:"like_count"`
	Visibility    string   `json:"visibility"`
	DownloadCount int      `json:"download_count"`
}

func toModelDTO(model *domain.Model) ModelDTO {
	dto := ModelDTO{
		Id:            model.Id.Identity(),
		Name:          model.Name.MSDName(),
		Owner:         model.Owner.Account(),
		Labels:        model.Labels,
		License:       model.License.License(),
		CreatedAt:     utils.ToDate(model.CreatedAt),
		UpdatedAt:     utils.ToDate(model.UpdatedAt),
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

type CmdToListModels struct {
	// can't define Name as domain.ResourceName
	// because the Name can be subpart of the real resource name
	Name string

	// list the models of Owner
	Owner primitive.Account

	// list models which have the labels
	Labels []string

	SortType primitive.SortType

	Count        bool
	PageNum      int
	CountPerPage int
}

func (cmd *CmdToListModels) toOption() repository.ListOption {
	v := repository.ListOption{
		Name:         cmd.Name,
		Owner:        cmd.Owner,
		Labels:       cmd.Labels,
		Count:        cmd.Count,
		PageNum:      cmd.PageNum,
		CountPerPage: cmd.CountPerPage,
	}

	if cmd.SortType == nil {
		v.SortType = primitive.SortTypeRecentlyUpdated
	} else {
		v.SortType = cmd.SortType
	}

	return v
}
