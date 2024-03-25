package domain

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type ActivityType string

const (
	Create ActivityType = "Create"
	Update ActivityType = "Update"
	Like   ActivityType = "Like"
)

// Activity struct represents the user activity entity.
type Activity struct {
	Owner    primitive.Account
	Type     ActivityType
	Time     int64
	Resource Resource
}

// Resource struct represents the resource object targeted by user activities.
type Resource struct {
	Type  primitive.ObjType  // Resource type
	Index primitive.Identity // Resource index
}
