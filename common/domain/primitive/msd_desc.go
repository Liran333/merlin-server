package primitive

import (
	"errors"

	"github.com/openmerlin/merlin-server/utils"
)

// MSDDesc
type MSDDesc interface {
	MSDDesc() string
}

func NewMSDDesc(v string) (MSDDesc, error) {
	if v == "" {
		return msdDesc(v), nil
	}

	v = utils.XSSEscapeString(v)
	if utils.StrLen(v) > msdConfig.MaxDescLength {
		return nil, errors.New("invalid desc")
	}

	return msdDesc(v), nil
}

func CreateMSDDesc(v string) MSDDesc {
	return msdDesc(v)
}

type msdDesc string

func (r msdDesc) MSDDesc() string {
	return string(r)
}
