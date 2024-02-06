package domain

import "encoding/json"

// spaceappCreatedEvent
type spaceappCreatedEvent struct {
	SpaceId  string `json:"space_id"`
	CommitId string `json:"commit_id"`
}

func (e *spaceappCreatedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewSpaceAppCreatedEvent(app *SpaceApp) spaceappCreatedEvent {
	return spaceappCreatedEvent{
		SpaceId:  app.SpaceId.Identity(),
		CommitId: app.CommitId,
	}
}
