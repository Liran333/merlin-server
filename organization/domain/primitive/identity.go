package primitive

import "errors"

const (
	legalRepresentative = "法定代表人"
	authorizedPerson    = "被授权人"
)

type Identity interface {
	Identity() string
}

func ValidateIdentity(v string) bool {
	return v == legalRepresentative || v == authorizedPerson
}

func NewIdentity(v string) (Identity, error) {
	if !ValidateIdentity(v) {
		return nil, errors.New("invalid identity")
	}

	return identity(v), nil
}

func CreateIdentity(v string) Identity {
	return identity(v)
}

type identity string

func (i identity) Identity() string {
	return string(i)
}
