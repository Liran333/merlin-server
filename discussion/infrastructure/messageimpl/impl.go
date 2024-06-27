package messageimpl

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"

	"github.com/openmerlin/merlin-server/discussion/domain/message"
)

func NewMessageImpl(t Topics) *messageImpl {
	return &messageImpl{
		topics: t,
	}
}

type messageImpl struct {
	topics Topics
}

func (impl messageImpl) SendUpdateCommentCountEvent(msg message.EventMessage) error {
	body, err := msg.Message()
	if err != nil {
		return err
	}

	return kfklib.Publish(impl.topics.CommentEvent, nil, body)
}
