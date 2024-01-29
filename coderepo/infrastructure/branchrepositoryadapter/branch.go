package branchrepositoryadapter

import (
	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type dao interface {
	DB() *gorm.DB
	GetRecord(filter, result interface{}) error
	DeleteByPrimaryKey(row interface{}) error
	EqualQuery(field string) string
}
type branchAdapter struct {
	dao
}

func (adapter *branchAdapter) Add(branch *domain.Branch) error {
	do := toBranchDO(branch)
	v := adapter.DB().Create(&do)

	return v.Error
}

func (adapter *branchAdapter) Delete(id primitive.Identity) error {
	return adapter.DeleteByPrimaryKey(
		&branchDO{Id: id.Integer()},
	)
}

func (adapter *branchAdapter) FindByIndex(index *domain.BranchIndex) (domain.Branch, error) {
	do := branchDO{
		Owner:  index.Owner.Account(),
		Repo:   index.Repo.MSDName(),
		Branch: index.Branch.BranchName(),
	}
	if err := adapter.GetRecord(&do, &do); err != nil {
		return domain.Branch{}, err
	}

	return do.toBranch(), nil
}
