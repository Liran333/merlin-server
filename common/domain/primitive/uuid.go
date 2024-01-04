package primitive

import "github.com/google/uuid"

// UUID
type UUID = uuid.UUID

func NewUUID(v string) (UUID, error) {
	return uuid.Parse(v)
}

func CreateUUID() UUID {
	return uuid.New()
}
