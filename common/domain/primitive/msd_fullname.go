package primitive

import (
	"errors"

	"github.com/openmerlin/merlin-server/utils"
)

// Fullname
type MSDFullname interface {
	MSDFullname() string
}

func NewMSDFullname(v string) (MSDFullname, error) {
	if v == "" {
		return msdFullname(v), nil
	}

	v = utils.XSSEscapeString(v)
	if utils.StrLen(v) > msdConfig.MaxFullnameLength {
		return nil, errors.New("invalid desc")
	}

	return msdFullname(v), nil
}

func CreateMSDFullname(v string) MSDFullname {
	return msdFullname(v)
}

type msdFullname string

func (r msdFullname) MSDFullname() string {
	return string(r)
}
