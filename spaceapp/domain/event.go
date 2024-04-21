/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain space app and configuration for the app service.
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

// Message returns the JSON representation of the spaceappRestartEvent.
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

// spaceappPausedEvent
type spaceappPausedEvent struct {
	SpaceId  string `json:"space_id"`
	CommitId string `json:"commit_id"`
}

// Message returns the JSON representation of the spaceappPausedEvent.
func (e *spaceappPausedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewSpaceAppPauseEvent creates a spaceappPausedEvent instance with the given SpaceApp.
func NewSpaceAppPauseEvent(app *SpaceAppIndex) spaceappPausedEvent {
	return spaceappPausedEvent{
		SpaceId:  app.SpaceId.Identity(),
		CommitId: app.CommitId,
	}
}

// spaceappResumedEvent
type spaceappResumedEvent struct {
	SpaceId  string `json:"space_id"`
	CommitId string `json:"commit_id"`
}

// Message returns the JSON representation of the spaceappResumedEvent.
func (e *spaceappResumedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewSpaceAppResumeEvent creates a spaceappResumedEvent instance with the given SpaceApp.
func NewSpaceAppResumeEvent(app *SpaceAppIndex) spaceappResumedEvent {
	return spaceappResumedEvent{
		SpaceId:  app.SpaceId.Identity(),
		CommitId: app.CommitId,
	}
}
