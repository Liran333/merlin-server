/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package spacerepositoryadapter

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/space/domain/repository"
)

const (
	// README represents the constant value for the README file.
	README = "README"
	readme = "readme"
)

type spaceAdapter struct {
	daoImpl
}

// Add adds a new space to the database and returns an error if any occurs.
func (adapter *spaceAdapter) Add(space *domain.Space) error {
	do := toSpaceDO(space)

	v := adapter.db().Create(&do)

	return v.Error
}

// FindByName finds a space by its name and returns it along with an error if any occurs.
func (adapter *spaceAdapter) FindByName(index *domain.SpaceIndex) (domain.Space, error) {
	do := spaceDO{Owner: index.Owner.Account(), Name: index.Name.MSDName()}

	if err := adapter.GetRecord(&do, &do); err != nil {
		return domain.Space{}, err
	}

	return do.toSpace(), nil
}

// FindById finds a space by its ID and returns it along with an error if any occurs.
func (adapter *spaceAdapter) FindById(spaceId primitive.Identity) (domain.Space, error) {
	do := spaceDO{Id: spaceId.Integer()}

	if err := adapter.GetByPrimaryKey(&do); err != nil {
		return domain.Space{}, err
	}

	return do.toSpace(), nil
}

// Delete deletes a space from the database by its ID and returns an error if any occurs.
func (adapter *spaceAdapter) Delete(spaceId primitive.Identity) error {
	return adapter.DeleteByPrimaryKey(
		&spaceDO{Id: spaceId.Integer()},
	)
}

// Save updates a space in the database and returns an error if any occurs.
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

// List is a method of spaceAdapter that takes a ListOption pointer as input
// and returns a slice of SpaceSummary, total count, and an error if any occurs.
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

// Count is a method of spaceAdapter that takes a ListOption pointer as input
// and returns the total count of spaces and an error if any occurs.
func (adapter *spaceAdapter) Count(opt *repository.ListOption) (int, error) {
	var total int64
	err := adapter.toQuery(opt).Count(&total).Error

	return int(total), err
}

func (adapter *spaceAdapter) toQuery(opt *repository.ListOption) *gorm.DB {
	db := adapter.db()

	if opt.Name != "" {
		query1, arg1 := likeFilter(fieldName, opt.Name)
		query2, arg2 := likeFilter(fieldFullName, opt.Name)
		db = db.Where(db.Where(query1, arg1).Or(query2, arg2))

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

func (adapter *spaceAdapter) SearchSpace(opt *repository.ListOption, login primitive.Account) ([]repository.SpaceSummary, int, error) {
	db := adapter.db()

	queryName, argName := likeFilter(fieldName, opt.Name)

	db = db.Where(queryName, argName)

	if login != nil {
		sql := fmt.Sprintf(`%s = ? or %s = ?`, fieldVisibility, fieldOwner)
		db = db.Where(sql, primitive.VisibilityPublic, login)
	} else {
		db = db.Where(equalQuery(fieldVisibility), primitive.VisibilityPublic)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if b, offset := opt.Pagination(); b {
		if offset > 0 {
			db = db.Limit(opt.CountPerPage).Offset(offset)
		} else {
			db = db.Limit(opt.CountPerPage)
		}
	}

	var dos []spaceDO

	if err := db.Find(&dos).Error; err != nil {
		return nil, 0, err
	}

	r := make([]repository.SpaceSummary, len(dos))
	for i := range dos {
		r[i] = dos[i].toSpaceSummary()
	}

	return r, int(total), nil
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

	default:
		return ""
	}
}

func (adapter *spaceAdapter) AddLike(spaceId primitive.Identity) error {
	id := spaceId.Integer()

	return adapter.IncrementLikeCount(id)
}

func (adapter *spaceAdapter) DeleteLike(spaceId primitive.Identity) error {
	id := spaceId.Integer()

	return adapter.DescendLikeCount(id)
}
