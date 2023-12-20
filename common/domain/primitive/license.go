package primitive

import (
	"errors"
	"strings"
)

// License
type License interface {
	License() string
}

func NewLicense(v string) (License, error) {
	v = strings.ToLower(strings.TrimSpace(v))

	if v == "" || !allLicenses[v] {
		return nil, errors.New("unsupported license")
	}

	return license(v), nil
}

func CreateLicense(v string) License {
	return license(v)
}

type license string

func (r license) License() string {
	return string(r)
}
