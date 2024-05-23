/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package message provides interfaces for defining event messages and sending dataset-related events.
package message

// EventMessage is an interface that represents an event message.
type EventMessage interface {
	Message() ([]byte, error)
}

// DatasetMessage is an interface that defines methods for sending dataset-related events.
type DatasetMessage interface {
	SendDatasetCreatedEvent(EventMessage) error
	SendDatasetDeletedEvent(EventMessage) error
	SendDatasetUpdatedEvent(EventMessage) error
}
