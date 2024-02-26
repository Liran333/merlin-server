/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package primitive

import (
	"errors"

	"github.com/openmerlin/merlin-server/utils"
)

// MSDFullname is an interface representing a full name.
type MSDFullname interface {
	MSDFullname() string
}

// NewMSDFullname creates a new MSDFullname instance from a string value.
func NewMSDFullname(v string) (MSDFullname, error) {
	if v == "" {
		return msdFullname(v), nil
	}

	v = utils.XSSEscapeString(v)
	if utils.StrLen(v) > msdConfig.MaxFullnameLength {
		return nil, errors.New("invalid fullname")
	}

	return msdFullname(v), nil
}

// CreateMSDFullname creates a new MSDFullname instance directly from a string value.
func CreateMSDFullname(v string) MSDFullname {
	return msdFullname(v)
}

type msdFullname string

// MSDFullname returns the string representation of the full name.
func (r msdFullname) MSDFullname() string {
	return string(r)
}
