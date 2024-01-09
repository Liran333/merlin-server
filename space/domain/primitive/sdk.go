package primitive

import (
	"errors"
	"strings"
)

type SDK interface {
	SDK() string
}

func NewSDK(v string) (SDK, error) {
	v = strings.ToLower(strings.TrimSpace(v))

	if v == "" || !allSDK[v] {
		return nil, errors.New("unsupported sdk")
	}

	return sdk(v), nil
}

func CreateSDK(v string) SDK {
	return sdk(v)
}

type sdk string

func (r sdk) SDK() string {
	return string(r)
}
