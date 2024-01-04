package primitive

import (
	"errors"
	"strings"
)

const (
	merlin = "merlin"
)

// UserAgent
type UserAgent interface {
	UserAgent() string
}

func NewUserAgent(v string) (UserAgent, error) {
	v = strings.ToLower(v)

	if v != merlin {
		return nil, errors.New("unknown user agent")
	}

	return userAgent(v), nil
}

func CreateUserAgent(v string) UserAgent {
	return userAgent(v)
}

type userAgent string

func (r userAgent) UserAgent() string {
	return string(r)
}
