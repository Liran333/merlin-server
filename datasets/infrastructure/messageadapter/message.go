/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package messageadapter provides an adapter for working with message-related functionality.
package messageadapter

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"

	"github.com/openmerlin/merlin-server/datasets/domain/message"
)

// MessageAdapter is a function that takes a pointer to a Topics struct and returns a pointer to a messageAdapter struct
// with the topics field set to the value of the input pointer.
func MessageAdapter(topic *Topics) *messageAdapter {
	return &messageAdapter{topics: *topic}
}

type messageAdapter struct {
	topics Topics
}

// SendDatasetCreatedEvent is a method on the messageAdapter struct that takes an EventMessage
// and sends it to the DatasetCreated topic.
func (p *messageAdapter) SendDatasetCreatedEvent(e message.EventMessage) error {
	return send(p.topics.DatasetCreated, e)
}

// SendDatasetDeletedEvent is a method on the messageAdapter struct that takes an EventMessage
// and sends it to the DatasetDeleted topic.
func (p *messageAdapter) SendDatasetDeletedEvent(e message.EventMessage) error {
	return send(p.topics.DatasetDeleted, e)
}

// SendDatasetUpdatedEvent is a method on the messageAdapter struct that takes an EventMessage
// and sends it to the DatasetUpdated topic.
func (p *messageAdapter) SendDatasetUpdatedEvent(e message.EventMessage) error {
	return send(p.topics.DatasetUpdated, e)
}

func send(topic string, v message.EventMessage) error {
	body, err := v.Message()
	if err != nil {
		return err
	}

	return kfklib.Publish(topic, nil, body)
}
