package primitive

import (
	"errors"
	"strings"
)

type Hardware interface {
	Hardware() string
}

func NewHardware(v string) (Hardware, error) {
	v = strings.ToLower(strings.TrimSpace(v))

	if v == "" || !allHardware[v] {
		return nil, errors.New("unsupported hardware")
	}

	return hardware(v), nil
}

func CreateHardware(v string) Hardware {
	return hardware(v)
}

type hardware string

func (r hardware) Hardware() string {
	return string(r)
}
