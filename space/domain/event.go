/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package domain

import (
	"encoding/json"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// spaceCreatedEvent
type spaceCreatedEvent struct {
	Time      int64  `json:"time"`
	Owner     string `json:"owner"`
	SpaceId   string `json:"space_id"`
	CreatedBy string `json:"created_by"`
}

// Message serializes the spaceCreatedEvent into a JSON byte array.
func (e *spaceCreatedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewSpaceCreatedEvent return a spaceCreatedEvent
func NewSpaceCreatedEvent(space *Space) spaceCreatedEvent {
	return spaceCreatedEvent{
		Time:      space.CreatedAt,
		Owner:     space.Owner.Account(),
		SpaceId:   space.Id.Identity(),
		CreatedBy: space.CreatedBy.Account(),
	}
}

// spaceDeletedEvent
type spaceDeletedEvent struct {
	SpaceId   string `json:"space_id"`
	DeletedBy string `json:"deleted_by"`
}

// Message serializes the spaceDeletedEvent into a JSON byte array.
func (e *spaceDeletedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewSpaceDeletedEvent creates a new spaceDeletedEvent instance with the given Space.
func NewSpaceDeletedEvent(user primitive.Account, space *Space) spaceDeletedEvent {
	return spaceDeletedEvent{
		SpaceId:   space.Id.Identity(),
		DeletedBy: user.Account(),
	}
}

// spaceUpdatedEvent
type spaceUpdatedEvent struct {
	Time       int64  `json:"time"`
	Repo       string `json:"repo"`
	Owner      string `json:"owner"`
	SpaceId    string `json:"space_id"`
	UpdatedBy  string `json:"updated_by"`
	IsPriToPub bool   `json:"is_pri_to_pub"` // private to public
}

// Message serializes the spaceUpdatedEvent into a JSON byte array.
func (e *spaceUpdatedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewSpaceUpdatedEvent return a spaceUpdatedEvent
func NewSpaceUpdatedEvent(user primitive.Account, space *Space, b bool) spaceUpdatedEvent {
	return spaceUpdatedEvent{
		Time:       space.UpdatedAt,
		Repo:       space.Name.MSDName(),
		Owner:      space.Owner.Account(),
		SpaceId:    space.Id.Identity(),
		UpdatedBy:  user.Account(),
		IsPriToPub: b,
	}
}

// spaceEnvChangedEvent
type spaceEnvChangedEvent struct {
	SpaceId   string `json:"space_id"`
	ChangedBy string `json:"changed_by"`
}

// Message serializes the spaceEnvChangedEvent into a JSON byte array.
func (e *spaceEnvChangedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewSpaceDeletedEvent creates a new spaceDeletedEvent instance with the given Space.
func NewSpaceEnvChangedEvent(user primitive.Account, space *Space) spaceEnvChangedEvent {
	return spaceEnvChangedEvent{
		SpaceId:   space.Id.Identity(),
		ChangedBy: user.Account(),
	}
}