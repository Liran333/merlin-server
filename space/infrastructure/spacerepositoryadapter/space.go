package spacerepositoryadapter

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/space/domain/repository"
)

const (
	README = "README"
	readme = "readme"
)

type spaceAdapter struct {
	daoImpl
}

func (adapter *spaceAdapter) Add(space *domain.Space) error {
	do := toSpaceDO(space)

	v := adapter.db().Create(&do)
	// TODO check if the space exists

	return v.Error
}

func (adapter *spaceAdapter) FindByName(index *domain.SpaceIndex) (domain.Space, error) {
	do := spaceDO{Owner: index.Owner.Account(), Name: index.Name.MSDName()}

	if err := adapter.GetRecord(&do, &do); err != nil {
		return domain.Space{}, err
	}

	return do.toSpace(), nil
}

func (adapter *spaceAdapter) FindById(spaceId primitive.Identity) (domain.Space, error) {
	do := spaceDO{Id: spaceId.Integer()}

	if err := adapter.GetByPrimaryKey(&do); err != nil {
		return domain.Space{}, err
	}

	return do.toSpace(), nil
}

func (adapter *spaceAdapter) Delete(spaceId primitive.Identity) error {
	return adapter.DeleteByPrimaryKey(
		&spaceDO{Id: spaceId.Integer()},
	)
}

func (adapter *spaceAdapter) Save(space *domain.Space) error {
	do := toSpaceDO(space)
	do.Version += 1

	v := adapter.db().Model(
		&spaceDO{Id: space.Id.Integer()},
	).Where(
		equalQuery(fieldVersion), space.Version,
	).Select(`*`).Updates(&do)

	if v.Error != nil {
		return v.Error
	}

	if v.RowsAffected == 0 {
		return commonrepo.NewErrorConcurrentUpdating(
			errors.New("concurrent updating"),
		)
	}

	return nil
}

func (adapter *spaceAdapter) List(opt *repository.ListOption) ([]repository.SpaceSummary, int, error) {
	query := adapter.toQuery(opt)

	// total
	var total int64
	if opt.Count {
		if err := query.Count(&total).Error; err != nil {
			return nil, 0, err
		}
	}

	// list
	if b, offset := opt.Pagination(); b {
		if offset > 0 {
			query = query.Limit(opt.CountPerPage).Offset(offset)
		} else {
			query = query.Limit(opt.CountPerPage)
		}
	}

	if v := order(opt.SortType); v != "" {
		query = query.Order(v)
	}

	var dos []spaceDO

	err := query.Find(&dos).Error
	if err != nil || len(dos) == 0 {
		return nil, 0, nil
	}

	r := make([]repository.SpaceSummary, len(dos))
	for i := range dos {
		r[i] = dos[i].toSpaceSummary()
	}

	return r, int(total), nil
}

func (adapter *spaceAdapter) toQuery(opt *repository.ListOption) *gorm.DB {
	db := adapter.db()

	if opt.Name != "" {
		query, arg := likeFilter(fieldName, opt.Name)

		db = db.Where(query, arg)

		if strings.Contains(readme, strings.ToLower(opt.Name)) {
			db = db.Where(notEqualQuery(fieldName), README)
		}
	} else {
		db = db.Where(notEqualQuery(fieldName), README)
	}

	if opt.Owner != nil {
		db = db.Where(equalQuery(fieldOwner), opt.Owner.Account())
	}

	if opt.Visibility != nil {
		db = db.Where(equalQuery(fieldVisibility), opt.Visibility.Visibility())
	}

	if opt.License != nil {
		db = db.Where(equalQuery(fieldLicense), opt.License.License())
	}

	if v := opt.Labels.Task; v != "" {
		db = db.Where(equalQuery(fieldTask), v)
	}

	if v := opt.Labels.Others; v != nil && v.Len() > 0 {
		query, arg := intersectionFilter(fieldOthers, v.UnsortedList())

		db = db.Where(query, arg)
	}

	if v := opt.Labels.Frameworks; v != nil && v.Len() > 0 {
		query, arg := intersectionFilter(fieldFrameworks, v.UnsortedList())

		db = db.Where(query, arg)
	}

	return db
}

func order(t primitive.SortType) string {
	if t == nil {
		return ""
	}

	switch t.SortType() {
	case primitive.SortByAlphabetical:
		return fieldName

	case primitive.SortByRecentlyUpdated:
		return orderByDesc(fieldUpdatedAt)

	case primitive.SortByRecentlyCreated:
		return orderByDesc(fieldCreatedAt)

	// TODO other type

	default:
		return ""
	}
}
