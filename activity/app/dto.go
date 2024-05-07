package app

import (
	"github.com/openmerlin/merlin-server/activity/domain"
	"github.com/openmerlin/merlin-server/activity/domain/repository"
	modelapp "github.com/openmerlin/merlin-server/models/app"
	spaceapp "github.com/openmerlin/merlin-server/space/app"
)

type ActivitysDTO struct {
	Total      int                  `json:"total"`
	AvatarId   string               `json:"avatar_id"`
	Activities []ActivitySummaryDTO `json:"activity"`
}

type ActivitySummaryDTO struct {
	ActivityDTO
	StatDTO
}

type ActivityDTO struct {
	Type     string      `json:"type"`
	Time     int64       `json:"time"`
	Name     string      `json:"name"`
	Owner    string      `json:"owner"`
	Resource ResourceDTO `json:"resource"`
	AdditionalDTO
}

type StatDTO struct {
	LikeCount     int `json:"like_count"`
	DownloadCount int `json:"download_count"`
}

type ResourceDTO struct {
	Type    string `json:"type"`
	Index   int64  `json:"index"`
	Owner   string `json:"owner"`
	Disable bool   `json:"disable"`
}

type AdditionalDTO struct {
	CustomerModelDTO
	CustomerSpaceDTO
	Fullname string `json:"fullname"`
}

type CustomerModelDTO struct {
	Labels modelapp.ModelLabelsDTO `json:"model_labels"`
}

type CustomerSpaceDTO struct {
	AvatarId string `json:"space_avatar_id"`
}

type CmdToAddActivity = domain.Activity
type CmdToListActivities = repository.ListOption

func toStatDTO(stat *domain.Stat) StatDTO {
	return StatDTO{
		LikeCount:     stat.LikeCount,
		DownloadCount: stat.DownloadCount,
	}
}

func toResourceDTO(resource *domain.Resource) ResourceDTO {
	return ResourceDTO{
		Type:    string(resource.Type),
		Index:   resource.Index.Integer(),
		Owner:   resource.Owner.Account(),
		Disable: resource.Disable,
	}
}

func toActivityDTO(activity *domain.Activity, additions AdditionalDTO) ActivityDTO {
	return ActivityDTO{
		Type:          string(activity.Type),
		Time:          activity.Time,
		Name:          activity.Name.MSDName(),
		Owner:         activity.Owner.Account(),
		Resource:      toResourceDTO(&activity.Resource),
		AdditionalDTO: additions,
	}
}

func fromModelDTO(model modelapp.ModelDTO, activity *domain.Activity, stat *domain.Stat) ActivitySummaryDTO {
	additions := AdditionalDTO{CustomerModelDTO: CustomerModelDTO{
		Labels: model.Labels,
	},
		Fullname: model.Fullname}

	return ActivitySummaryDTO{
		ActivityDTO: toActivityDTO(activity, additions),
		StatDTO:     toStatDTO(stat),
	}
}

func fromSpaceDTO(space spaceapp.SpaceDTO, activity *domain.Activity, stat *domain.Stat) ActivitySummaryDTO {
	additions := AdditionalDTO{CustomerSpaceDTO: CustomerSpaceDTO{
		AvatarId: space.AvatarId,
	},
		Fullname: space.Fullname,
	}

	return ActivitySummaryDTO{
		ActivityDTO: toActivityDTO(activity, additions),
		StatDTO:     toStatDTO(stat),
	}
}
