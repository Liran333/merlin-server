/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package messageadapter

import (
	"github.com/openmerlin/merlin-server/models/domain/message"
	kfklib "github.com/opensourceways/kafka-lib/agent"
)

// MessageAdapter is a function that takes a pointer to a Topics struct and returns a pointer to a messageAdapter struct
// with the topics field set to the value of the input pointer.
func MessageAdapter(topic *Topics) *messageAdapter {
	return &messageAdapter{topics: *topic}
}

type messageAdapter struct {
	topics Topics
}

// SendModelCreatedEvent is a method on the messageAdapter struct that takes an EventMessage
// and sends it to the ModelCreate topic.
func (p *messageAdapter) SendModelCreatedEvent(e message.EventMessage) error {
	return send(p.topics.ModelCreated, e)
}

// SendModelDeletedEvent is a method on the messageAdapter struct that takes an EventMessage
// and sends it to the ModelDeleted topic.
func (p *messageAdapter) SendModelDeletedEvent(e message.EventMessage) error {
	return send(p.topics.ModelDeleted, e)
}

// SendModelUpdatedEvent is a method on the messageAdapter struct that takes an EventMessage
// and sends it to the ModelUpdated topic.
func (p *messageAdapter) SendModelUpdatedEvent(e message.EventMessage) error {
	return send(p.topics.ModelUpdated, e)
}

func send(topic string, v message.EventMessage) error {
	body, err := v.Message()
	if err != nil {
		return err
	}

	return kfklib.Publish(topic, nil, body)
}
