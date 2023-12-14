package primitive

import (
	"errors"
	"regexp"
)

var regMSDName = regexp.MustCompile("^[a-zA-Z0-9_-]+$")

// Name
type MSDName interface {
	MSDName() string
	FirstLetter() byte
}

func NewMSDName(v string) (MSDName, error) {
	n := len(v)
	if n > msdConfig.MaxNameLength || n < msdConfig.MinNameLength {
		return nil, errors.New("invalid name")
	}

	if !regMSDName.MatchString(v) {
		return nil, errors.New("invalid name")
	}

	return msdName(v), nil
}

func CreateMSDName(v string) MSDName {
	return msdName(v)
}

type msdName string

func (r msdName) MSDName() string {
	return string(r)
}

func (r msdName) FirstLetter() byte {
	return string(r)[0]
}
