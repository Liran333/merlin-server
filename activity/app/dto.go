package app

import (
	"github.com/openmerlin/merlin-server/activity/domain"
	"github.com/openmerlin/merlin-server/activity/domain/repository"
)

// AcctivityDTO is a struct that represents a data transfer object for an activity.
type AcctivityDTO struct {
	Total      int               `json:"total"`
	Activities []domain.Activity `json:"activity"`
}

type CmdToAddActivity = domain.Activity

// CmdToListActivities is a type alias for repository.ListOption, representing a command to list models.
type CmdToListActivities = repository.ListOption
