package app

import (
	"github.com/openmerlin/merlin-server/activity/domain"
	"github.com/openmerlin/merlin-server/activity/domain/repository"
	datasetapp "github.com/openmerlin/merlin-server/datasets/app"
	modelapp "github.com/openmerlin/merlin-server/models/app"
	spaceapp "github.com/openmerlin/merlin-server/space/app"
)

// ActivitysDTO represents the DTO (Data Transfer Object) structure for activities.
type ActivitysDTO struct {
	Total      int                  `json:"total"`
	AvatarId   string               `json:"avatar_id"`
	Activities []ActivitySummaryDTO `json:"activity"`
}

// ActivitySummaryDTO represents the DTO structure for activity summaries.
type ActivitySummaryDTO struct {
	ActivityDTO
	StatDTO
}

// ActivityDTO represents the DTO structure for an activity.
type ActivityDTO struct {
	Type     string      `json:"type"`
	Time     int64       `json:"time"`
	Name     string      `json:"name"`
	Owner    string      `json:"owner"`
	Resource ResourceDTO `json:"resource"`
	AdditionalDTO
}

// StatDTO represents the DTO structure for statistics.
type StatDTO struct {
	LikeCount     int `json:"like_count"`
	DownloadCount int `json:"download_count"`
}

// ResourceDTO represents the DTO structure for a resource.
type ResourceDTO struct {
	Type    string `json:"type"`
	Index   int64  `json:"index"`
	Owner   string `json:"owner"`
	Disable bool   `json:"disable"`
}

// AdditionalDTO represents the DTO structure for additional data.
type AdditionalDTO struct {
	CustomerModelDTO
	CustomerDatasetDTO
	CustomerSpaceDTO
	Fullname string `json:"fullname"`
}

// CustomerModelDTO represents the DTO (Data Transfer Object) structure for customer models.
type CustomerModelDTO struct {
	Labels modelapp.ModelLabelsDTO `json:"model_labels"`
}

// CustomerDatasetDTO represents the DTO (Data Transfer Object) structure for customer datasets.
type CustomerDatasetDTO struct {
	Labels datasetapp.DatasetLabelsDTO `json:"dataset_labels"`
}

// CustomerSpaceDTO represents the DTO structure for customer spaces.
type CustomerSpaceDTO struct {
	AvatarId string `json:"space_avatar_id"`
}

// CmdToAddActivity is an alias for the domain.Activity type.
type CmdToAddActivity = domain.Activity

// CmdToListActivities is an alias for the repository.ListOption type.
type CmdToListActivities = repository.ListOption

// toStatDTO converts a domain.Stat object to a StatDTO object.
func toStatDTO(stat *domain.Stat) StatDTO {
	return StatDTO{
		LikeCount:     stat.LikeCount,
		DownloadCount: stat.DownloadCount,
	}
}

// toResourceDTO converts a domain.Resource object to a ResourceDTO object.
func toResourceDTO(resource *domain.Resource) ResourceDTO {
	return ResourceDTO{
		Type:    string(resource.Type),
		Index:   resource.Index.Integer(),
		Owner:   resource.Owner.Account(),
		Disable: resource.Disable,
	}
}

// toActivityDTO converts a domain.Activity object to an ActivityDTO object.
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

// fromModelDTO converts a modelapp.ModelDTO, domain.Activity, and domain.Stat objects to an ActivitySummaryDTO object.
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

// fromDatasetDTO converts a datasetapp.DatasetDTO, domain.Activity, and domain.Stat objects to an ActivitySummaryDTO object.
func fromDatasetDTO(dataset datasetapp.DatasetDTO, activity *domain.Activity, stat *domain.Stat) ActivitySummaryDTO {
	additions := AdditionalDTO{CustomerDatasetDTO: CustomerDatasetDTO{
		Labels: dataset.Labels,
	},
		Fullname: dataset.Fullname}

	return ActivitySummaryDTO{
		ActivityDTO: toActivityDTO(activity, additions),
		StatDTO:     toStatDTO(stat),
	}
}

// fromSpaceDTO converts a spaceapp.SpaceDTO, domain.Activity, and domain.Stat objects to an ActivitySummaryDTO object.
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
