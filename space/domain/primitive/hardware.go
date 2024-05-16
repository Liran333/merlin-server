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
	IsNpu() bool
	IsCpu() bool
}

// NewHardware creates a new Hardware instance decided by sdk based on the given string.
func NewHardware(v string, sdk string) (Hardware, error) {
	v = strings.ToLower(strings.TrimSpace(v))
	sdk = strings.ToLower(strings.TrimSpace(sdk))

	if _, ok := sdkObjects[sdk]; sdk == "" || !ok {
		return nil, errors.New("unsupported sdk")
	}

	if v == "" || !sdkObjects[sdk].Has(v) {
		return nil, errors.New("unsupported hardware")
	}

	return hardware(v), nil
}

func IsValidHardware(h string) bool {
	for _, sdk := range sdkObjects {
		if sdk.Has(h) {
			return true
		}
	}

	return false
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

func (r hardware) IsNpu() bool {
	return strings.Contains(strings.ToLower(string(r)), "npu")
}

func (r hardware) IsCpu() bool {
	return strings.Contains(strings.ToLower(string(r)), "cpu")
}
