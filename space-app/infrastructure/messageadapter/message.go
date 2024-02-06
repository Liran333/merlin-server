package messageadapter

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"

	"github.com/openmerlin/merlin-server/space-app/domain/message"
)

func MessageAdapter(topic *Topics) *messageAdapter {
	return &messageAdapter{*topic}
}

type messageAdapter struct {
	topics Topics
}

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
