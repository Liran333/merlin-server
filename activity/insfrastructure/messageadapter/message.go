/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package messageadapter provides an adapter for working with message-related functionality.
package messageadapter

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	"golang.org/x/xerrors"

	"github.com/openmerlin/merlin-server/activity/domain/message"
)

// MessageAdapter is a function that takes a pointer to a Topics struct and returns a pointer to a messageAdapter struct
// with the topics field set to the value of the input pointer.
func MessageAdapter(topic *Topics) *messageAdapter {
	return &messageAdapter{topics: *topic}
}

type messageAdapter struct {
	topics Topics
}

// SendLikeCreatedEvent sends a LikeCreatedEvent message using the messageAdapter.
func (p *messageAdapter) SendLikeCreatedEvent(e message.EventMessage) error {
	err := send(p.topics.LikeCreate, e)
	if err != nil {
		return xerrors.Errorf("failed to send kafka message, error:%w", err)
	}
	return nil
}

// SendLikeDeletedEvent sends a LikeDeletedEvent message using the messageAdapter.
func (p *messageAdapter) SendLikeDeletedEvent(e message.EventMessage) error {
	return send(p.topics.LikeDelete, e)
}

func send(topic string, v message.EventMessage) error {
	body, err := v.Message()
	if err != nil {
		return xerrors.Errorf("failed to send like delete/create event, error:%w", err)
	}

	return kfklib.Publish(topic, nil, body)
}
