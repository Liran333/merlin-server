package primitive

import (
	"errors"
	"net/url"
)

// AvatarId
type AvatarId interface {
	AvatarId() string
}

func NewAvatarId(v string) (AvatarId, error) {
	if v == "" {
		return dpAvatarId(v), nil
	}

	avatarId, err := url.ParseRequestURI(v)
	if err != nil {
		return nil, errors.New("avatar must be a valid uri")
	}

	return dpAvatarId(avatarId.String()), nil
}

func CreateAvatarId(v string) AvatarId {
	return dpAvatarId(v)
}

type dpAvatarId string

func (r dpAvatarId) AvatarId() string {
	return string(r)
}

func (r dpAvatarId) DomainValue() string {
	return string(r)
}
