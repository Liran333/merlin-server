package message

type EventMessage interface {
	Message() ([]byte, error)
}

type SpaceAppMessage interface {
	SendSpaceAppCreatedEvent(EventMessage) error
}
