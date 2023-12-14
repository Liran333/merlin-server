package domain

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type CodeRepo struct {
	Id         string
	Name       primitive.MSDName
	Owner      primitive.Account
	License    primitive.License
	Visibility primitive.Visibility
}

func (r *CodeRepo) IsPrivate() bool {
	return r.Visibility.IsPrivate()
}

func (r *CodeRepo) IsPublic() bool {
	return r.Visibility.IsPublic()
}
