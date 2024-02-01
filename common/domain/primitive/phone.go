package primitive

import (
	"errors"
	"regexp"
)

var phoneRegexp = regexp.MustCompile(`^(\\+)?[0-9]+$`)

// Phone number, china mainland supported only for now
type Phone interface {
	PhoneNumber() string
}

func NewPhone(v string) (Phone, error) {
	if v == "" {
		return nil, errors.New("empty phone number")
	}

	if !phoneRegexp.MatchString(v) {
		return nil, errors.New("invalid name")
	}

	return phoneNumber(v), nil
}

func CreatePhoneNumber(v string) Phone {
	return phoneNumber(v)
}

type phoneNumber string

func (r phoneNumber) PhoneNumber() string {
	return string(r)
}
