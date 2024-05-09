/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package primitive

import (
	"golang.org/x/xerrors"
)

// ENVName is an interface that defines the method to get the ENV Name.
type ENVName interface {
	ENVName() string
}

// NewENVName creates a new ENV instance string name.
func NewENVName(v string) (ENVName, error) {
	n := len(v)
	if n > envConfig.MaxValueLength || n < envConfig.MinValueLength {
		return nil, xerrors.Errorf("invalid Value length, should between %d and %d",
			envConfig.MinValueLength, envConfig.MaxValueLength)
	}

	if !envConfig.nameRegexp.MatchString(v) {
		return nil, xerrors.Errorf("invalid name")
	}

	return envName(v), nil
}

// CreateENVName creates a new ENV instance from a string name.
func CreateENVName(v string) ENVName {
	return envName(v)
}

type envName string

// ENVName returns the string representation of the env name.
func (r envName) ENVName() string {
	return string(r)
}
