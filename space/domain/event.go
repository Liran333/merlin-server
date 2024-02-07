package domain

import "encoding/json"

// spaceDeletedEvent
type spaceDeletedEvent struct {
	SpaceId string `json:"space_id"`
}

func (e *spaceDeletedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewSpaceDeletedEvent(space *Space) spaceDeletedEvent {
	return spaceDeletedEvent{
		SpaceId: space.Id.Identity(),
	}
}

var NewSpaceUpdatedEvent = NewSpaceDeletedEvent
