/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package primitive

import (
	"errors"
	"strings"
)

// Hardware is an interface that defines hardware-related operations.
type Hardware interface {
	Hardware() string
}

// NewHardware creates a new Hardware instance decided by sdk based on the given string.
func NewHardware(v string, sdk string) (Hardware, error) {
	v = strings.ToLower(strings.TrimSpace(v))

	if _, ok := sdkObjects[sdk]; sdk == "" || !ok {
		return nil, errors.New("unsupported sdk")
	}

	if v == "" || !sdkObjects[sdk].Has(v) {
		return nil, errors.New("unsupported hardware")
	}

	return hardware(v), nil
}

// CreateHardware creates a new Hardware instance based on the given string.
func CreateHardware(v string) Hardware {
	return hardware(v)
}

type hardware string

// Hardware returns the string representation of the hardware.
func (r hardware) Hardware() string {
	return string(r)
}
