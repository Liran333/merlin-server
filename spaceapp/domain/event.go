/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package domain

import "encoding/json"

// spaceappCreatedEvent
type spaceappCreatedEvent struct {
	SpaceId  string `json:"space_id"`
	CommitId string `json:"commit_id"`
}

// Message returns the JSON representation of the spaceappCreatedEvent.
func (e *spaceappCreatedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewSpaceAppCreatedEvent creates a new spaceappCreatedEvent instance with the given SpaceApp.
func NewSpaceAppCreatedEvent(app *SpaceApp) spaceappCreatedEvent {
	return spaceappCreatedEvent{
		SpaceId:  app.SpaceId.Identity(),
		CommitId: app.CommitId,
	}
}

// spaceappRestartEvent
type spaceappRestartEvent struct {
	SpaceId  string `json:"space_id"`
	CommitId string `json:"commit_id"`
}

// Message returns the JSON representation of the spaceappCreatedEvent.
func (e *spaceappRestartEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewSpaceAppRestartEvent creates a spaceappRestartEvent instance with the given SpaceApp.
func NewSpaceAppRestartEvent(app *SpaceAppIndex) spaceappRestartEvent {
	return spaceappRestartEvent{
		SpaceId:  app.SpaceId.Identity(),
		CommitId: app.CommitId,
	}
}
