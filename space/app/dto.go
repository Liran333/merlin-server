package app

import (
	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/space/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

type CmdToCreateSpace struct {
	coderepoapp.CmdToCreateRepo

	Desc     primitive.MSDDesc
	Fullname primitive.MSDFullname
}

type CmdToUpdateSpace struct {
	coderepoapp.CmdToUpdateRepo

	Desc     primitive.MSDDesc
	Fullname primitive.MSDFullname
}

func (cmd *CmdToUpdateSpace) toSpace(space *domain.Space) (b bool) {
	if v := cmd.Desc; v != nil && v != space.Desc {
		space.Desc = v
		b = true
	}

	if v := cmd.Fullname; v != nil && v != space.Fullname {
		space.Fullname = v
		b = true
	}

	return
}

type SpaceDTO struct {
	Id            string         `json:"id"`
	Name          string         `json:"name"`
	Desc          string         `json:"desc"`
	Owner         string         `json:"owner"`
	Labels        SpaceLabelsDTO `json:"labels"`
	Fullname      string         `json:"fullname"`
	CreatedAt     string         `json:"created_at"`
	UpdatedAt     string         `json:"updated_at"`
	LikeCount     int            `json:"like_count"`
	Visibility    string         `json:"visibility"`
	DownloadCount int            `json:"download_count"`
}

type SpaceLabelsDTO struct {
	Task       string   `json:"task"`
	Others     []string `json:"others"`
	License    string   `json:"license"`
	Frameworks []string `json:"frameworks"`
}

func toSpaceLabelsDTO(space *domain.Space) SpaceLabelsDTO {
	labels := &space.Labels

	return SpaceLabelsDTO{
		Task:       labels.Task,
		Others:     labels.Others.UnsortedList(),
		License:    space.License.License(),
		Frameworks: labels.Frameworks.UnsortedList(),
	}
}

func toSpaceDTO(space *domain.Space) SpaceDTO {
	dto := SpaceDTO{
		Id:            space.Id.Identity(),
		Name:          space.Name.MSDName(),
		Owner:         space.Owner.Account(),
		Labels:        toSpaceLabelsDTO(space),
		CreatedAt:     utils.ToDate(space.CreatedAt),
		UpdatedAt:     utils.ToDate(space.UpdatedAt),
		LikeCount:     space.LikeCount,
		Visibility:    space.Visibility.Visibility(),
		DownloadCount: space.DownloadCount,
	}

	if space.Desc != nil {
		dto.Desc = space.Desc.MSDDesc()
	}

	if space.Fullname != nil {
		dto.Fullname = space.Fullname.MSDFullname()
	}

	return dto
}

type SpacesDTO struct {
	Total  int                       `json:"total"`
	Spaces []repository.SpaceSummary `json:"spaces"`
}

type CmdToListSpaces = repository.ListOption
