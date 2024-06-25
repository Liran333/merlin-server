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
	SpaceId string `json:"space_id"`
}

// Message returns the JSON representation of the spaceappPausedEvent.
func (e *spaceappPausedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewSpaceAppPauseEvent creates a spaceappPausedEvent instance with the given SpaceApp.
func NewSpaceAppPauseEvent(app *SpaceAppIndex) spaceappPausedEvent {
	return spaceappPausedEvent{
		SpaceId: app.SpaceId.Identity(),
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

// spaceappHeartbeatEvent
type spaceappHeartbeatEvent struct {
	SpaceId  	string `json:"space_id"`
	CommitId 	string `json:"commit_id"`
}

// Message returns the JSON representation of the spaceappResumedEvent.
func (e *spaceappHeartbeatEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewSpaceAppHeartbeatEvent creates a spaceappHeartbeatEvent instance with the given SpaceApp.
func NewSpaceAppHeartbeatEvent(app *SpaceApp) spaceappHeartbeatEvent {
	return spaceappHeartbeatEvent{
		SpaceId:  	app.SpaceId.Identity(),
		CommitId: 	app.CommitId,
	}
}

// spaceappSleepEvent
type spaceappSleepEvent struct {
	SpaceId  	string `json:"space_id"`
}

// Message returns the JSON representation of the spaceappSleepEvent.
func (e *spaceappSleepEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewSpaceAppSleepEvent creates a spaceappSleepEvent instance with the given SpaceApp.
func NewSpaceAppSleepEvent(app *SpaceAppIndex) spaceappSleepEvent {
	return spaceappSleepEvent{
		SpaceId:  	app.SpaceId.Identity(),
	}
}

// spaceappWakeupEvent
type spaceappWakeupEvent struct {
	SpaceId  	string `json:"space_id"`
}

// Message returns the JSON representation of the spaceappWakeupEvent.
func (e *spaceappWakeupEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewSpaceAppWakeupEvent creates a spaceappWakeupEvent instance with the given SpaceApp.
func NewSpaceAppWakeupEvent(app *SpaceAppIndex) spaceappWakeupEvent {
	return spaceappWakeupEvent{
		SpaceId:  	app.SpaceId.Identity(),
	}
}