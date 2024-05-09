/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package primitive

import (
	"golang.org/x/xerrors"
)

// ENVValue is an interface that defines the method to get the ENV Value.
type ENVValue interface {
	ENVValue() string
}

// NewENVValue creates a new ENV instance string value.
func NewENVValue(v string) (ENVValue, error) {
	n := len(v)
	if n > envConfig.MaxValueLength || n < envConfig.MinValueLength {
		return nil, xerrors.Errorf("invalid Value length, should between %d and %d",
			envConfig.MinValueLength, envConfig.MaxValueLength)
	}

	return envValue(v), nil
}

// CreateENVValue creates a new ENV instance from a string value.
func CreateENVValue(v string) ENVValue {
	return envValue(v)
}

type envValue string

// ENVValue returns the string representation of the env value.
func (r envValue) ENVValue() string {
	return string(r)
}
