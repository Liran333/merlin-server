/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package messageadapter

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"

	"github.com/openmerlin/merlin-server/spaceapp/domain/message"
)

// MessageAdapter creates a new messageAdapter instance with the given Topics.
func MessageAdapter(topic *Topics) *messageAdapter {
	return &messageAdapter{topics: *topic}
}

type messageAdapter struct {
	topics Topics
}

// SendSpaceAppCreatedEvent sends a SpaceAppCreated event message to the corresponding topic.
func (p *messageAdapter) SendSpaceAppCreatedEvent(e message.EventMessage) error {
	return send(p.topics.SpaceAppCreated, e)
}

func send(topic string, v message.EventMessage) error {
	body, err := v.Message()
	if err != nil {
		return err
	}

	return kfklib.Publish(topic, nil, body)
}
