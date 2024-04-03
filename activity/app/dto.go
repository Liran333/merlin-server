package app

import (
	"github.com/openmerlin/merlin-server/activity/domain"
	"github.com/openmerlin/merlin-server/activity/domain/repository"
)

// ActivityDTO is a struct that represents a data transfer object for an activity.
type ActivityDTO struct {
	Total      int               `json:"total"`
	AvatarId   string            `json:"avatar_id"`
	Activities []domain.Activity `json:"activity"`
}

type CmdToAddActivity = domain.Activity

// CmdToListActivities is a type alias for repository.ListOption, representing a command to list models.
type CmdToListActivities = repository.ListOption
