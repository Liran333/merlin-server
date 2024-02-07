package message

type EventMessage interface {
	Message() ([]byte, error)
}

type SpaceMessage interface {
	SendSpaceDeletedEvent(EventMessage) error
	SendSpaceUpdatedEvent(EventMessage) error
}
