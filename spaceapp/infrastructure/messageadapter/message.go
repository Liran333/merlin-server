/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package messageadapter

import (
	"fmt"

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

// SendSpaceAppRestartedEvent sends a SpaceAppRestarted event message to the corresponding topic.
func (p *messageAdapter) SendSpaceAppRestartedEvent(e message.EventMessage) error {
	return send(p.topics.SpaceAppRestarted, e)
}

func send(topic string, v message.EventMessage) error {
	body, err := v.Message()
	if err != nil {
		return fmt.Errorf("send msg topic:%s err:%w", topic, err)
	}

	err = kfklib.Publish(topic, nil, body)
	if err != nil {
		err = fmt.Errorf("send publish topic:%s err:%w", topic, err)
	}
	return err
}
