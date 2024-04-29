/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package message provides functionality for sending and handling event messages.
package message

// EventMessage is an interface that represents an event message.
type EventMessage interface {
	Message() ([]byte, error)
}

// SpaceAppMessage is an interface that defines a method for sending a space app created event.
type SpaceAppMessage interface {
	SendSpaceAppCreatedEvent(EventMessage) error
	SendSpaceAppRestartedEvent(EventMessage) error
	SendSpaceAppPauseEvent(EventMessage) error
	SendSpaceAppResumeEvent(EventMessage) error
	SendSpaceAppForcePauseEvent(EventMessage) error
}
