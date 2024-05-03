/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package messageadapter

import (
	"fmt"

	kfklib "github.com/opensourceways/kafka-lib/agent"

	"github.com/openmerlin/merlin-server/computility/domain/message"
)

// MessageAdapter creates a new messageAdapter instance with the given Topics.
func MessageAdapter(topic *Topics) *messageAdapter {
	return &messageAdapter{topics: *topic}
}

type messageAdapter struct {
	topics Topics
}

// SendComputilityRecallEvent sends a reall quota event message to the corresponding topic.
func (p *messageAdapter) SendComputilityRecallEvent(e message.EventMessage) error {
	return send(p.topics.ComputilityRecalled, e)
}

func send(topic string, v message.EventMessage) error {
	body, err := v.Message()
	if err != nil {
		return fmt.Errorf("send msg topic:%s err: %w", topic, err)
	}

	err = kfklib.Publish(topic, nil, body)
	if err != nil {
		err = fmt.Errorf("send publish topic:%s err: %w", topic, err)
	}

	return err
}
