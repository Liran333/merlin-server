package primitive

import (
	"errors"
	"strings"
)

const (
	Public  = "public"
	Private = "private"
)

var (
	VisibilityPublic  = visibility(Public)
	VisibilityPrivate = visibility(Private)
)

// Visibility
type Visibility interface {
	IsPublic() bool
	IsPrivate() bool
	Visibility() string
}

func NewVisibility(v string) (Visibility, error) {
	v = strings.ToLower(v)
	if v != Public && v != Private {
		return nil, errors.New("unknown visibility")
	}

	return visibility(v), nil
}

func CreateVisibility(v string) Visibility {
	return visibility(v)
}

type visibility string

func (r visibility) Visibility() string {
	return string(r)
}

func (r visibility) IsPrivate() bool {
	return string(r) == Private
}

func (r visibility) IsPublic() bool {
	return string(r) == Public
}
