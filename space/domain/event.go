/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package domain

import "encoding/json"

// spaceDeletedEvent
type spaceDeletedEvent struct {
	SpaceId string `json:"space_id"`
}

// Message serializes the spaceDeletedEvent into a JSON byte array.
func (e *spaceDeletedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewSpaceDeletedEvent creates a new spaceDeletedEvent instance with the given Space.
func NewSpaceDeletedEvent(space *Space) spaceDeletedEvent {
	return spaceDeletedEvent{
		SpaceId: space.Id.Identity(),
	}
}

// NewSpaceUpdatedEvent is an alias for NewSpaceDeletedEvent.
var NewSpaceUpdatedEvent = NewSpaceDeletedEvent
