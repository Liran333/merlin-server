/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package messageadapter

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"

	"github.com/openmerlin/merlin-server/space/domain/message"
)

// MessageAdapter is a function that takes a pointer to a Topics struct and returns a pointer to a messageAdapter struct
// with the topics field set to the value of the input pointer.
func MessageAdapter(topic *Topics) *messageAdapter {
	return &messageAdapter{topics: *topic}
}

type messageAdapter struct {
	topics Topics
}

// SendSpaceCreatedEvent is a method on the messageAdapter struct that takes an EventMessage
// and sends it to the SpaceCreated topic.
func (p *messageAdapter) SendSpaceCreatedEvent(e message.EventMessage) error {
	return send(p.topics.SpaceCreated, e)
}

// SendSpaceDeletedEvent is a method on the messageAdapter struct that takes an EventMessage
// and sends it to the SpaceDeleted topic.
func (p *messageAdapter) SendSpaceDeletedEvent(e message.EventMessage) error {
	return send(p.topics.SpaceDeleted, e)
}

// SendSpaceUpdatedEvent is a method on the messageAdapter struct that takes an EventMessage
// and sends it to the SpaceUpdated topic.
func (p *messageAdapter) SendSpaceUpdatedEvent(e message.EventMessage) error {
	return send(p.topics.SpaceUpdated, e)
}

// SendSpaceEnvChangedEvent is a method on the messageAdapter struct that takes an EventMessage
// and sends it to the SpaceEnvChanged topic.
func (p *messageAdapter) SendSpaceEnvChangedEvent(e message.EventMessage) error {
	return send(p.topics.SpaceEnvChanged, e)
}

func send(topic string, v message.EventMessage) error {
	body, err := v.Message()
	if err != nil {
		return err
	}

	return kfklib.Publish(topic, nil, body)
}
