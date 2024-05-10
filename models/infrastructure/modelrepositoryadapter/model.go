/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package modelrepositoryadapter provides an adapter for the model repository
package modelrepositoryadapter

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/models/domain"
	"github.com/openmerlin/merlin-server/models/domain/repository"
	orgrepo "github.com/openmerlin/merlin-server/organization/domain/repository"
)

const (
	// README represents the constant value for the README file.
	README = "README"
	readme = "readme"
)

type modelAdapter struct {
	daoImpl
}

// Add adds a new model to the database.
func (adapter *modelAdapter) Add(model *domain.Model) error {
	do := toModelDO(model)

	v := adapter.db().Create(&do)

	return v.Error
}

// FindByName finds a model by its name.
func (adapter *modelAdapter) FindByName(index *domain.ModelIndex) (domain.Model, error) {
	do := modelDO{Owner: index.Owner.Account(), Name: index.Name.MSDName()}

	if err := adapter.GetLowerModelName(&do, &do); err != nil {
		return domain.Model{}, err
	}

	return do.toModel(), nil
}

// FindById finds a model by its ID.
func (adapter *modelAdapter) FindById(modelId primitive.Identity) (domain.Model, error) {
	do := modelDO{Id: modelId.Integer()}

	if err := adapter.GetByPrimaryKey(&do); err != nil {
		return domain.Model{}, err
	}

	return do.toModel(), nil
}

// Delete deletes a model by its ID.
func (adapter *modelAdapter) Delete(modelId primitive.Identity) error {
	return adapter.DeleteByPrimaryKey(
		&modelDO{Id: modelId.Integer()},
	)
}

// Save updates an existing model in the database.
func (adapter *modelAdapter) Save(model *domain.Model) error {
	do := toModelDO(model)
	do.Version += 1

	v := adapter.db().Model(
		&modelDO{Id: model.Id.Integer()},
	).Where(
		equalQuery(fieldVersion), model.Version,
	).Select(`*`).Omit(fieldTask, fieldOthers, fieldFrameworks, filedLibraryName).Updates(&do)

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

// List retrieves a list of models based on the provided options.
func (adapter *modelAdapter) List(opt *repository.ListOption, login primitive.Account, member orgrepo.OrgMember) ([]repository.ModelSummary, int, error) {
	query := adapter.toQuery(opt)

	if opt.Visibility != nil {
		if login != nil {
			members, err := member.GetByUser(login.Account())
			if err != nil {
				return nil, 0, err
			}
			orgNames := make([]string, len(members))
			for _, member := range members {
				orgNames = append(orgNames, member.OrgName.Account())
			}
			sql := fmt.Sprintf(`%s = ? or %s = ? or %s in (?)`, fieldVisibility, fieldOwner, fieldOwner)
			query = query.Where(sql, opt.Visibility, login, orgNames)
		} else {
			query = query.Where(equalQuery(fieldVisibility), opt.Visibility.Visibility())
		}
	}

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

	var dos []modelDO

	err := query.Find(&dos).Error
	if err != nil || len(dos) == 0 {
		return nil, 0, nil
	}

	r := make([]repository.ModelSummary, len(dos))
	for i := range dos {
		r[i] = dos[i].toModelSummary()
	}

	return r, int(total), nil
}

// Count counts the number of models based on the provided options.
func (adapter *modelAdapter) Count(opt *repository.ListOption) (int, error) {
	var total int64
	err := adapter.toQuery(opt).Count(&total).Error

	return int(total), err
}

func (adapter *modelAdapter) toQuery(opt *repository.ListOption) *gorm.DB {
	db := adapter.db()

	if opt.Name != "" {
		_, arg := likeFilter(fieldName, opt.Name)
		db = db.Where(gorm.Expr("CONCAT("+fieldOwner+", '/', "+fieldName+") ilike ?", arg))

		if !opt.ExcludeFullname {
			query2, arg2 := likeFilter(fieldFullName, opt.Name)
			db = db.Or(query2, arg2)
		}

		if strings.Contains(readme, strings.ToLower(opt.Name)) {
			db = db.Where(notEqualQuery(fieldName), README)
		}
	} else {
		db = db.Where(notEqualQuery(fieldName), README)
	}

	if opt.Owner != nil {
		db = db.Where(equalQuery(fieldOwner), opt.Owner.Account())
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

	case primitive.SortByMostDownloads:
		return orderByDesc(fieldDownloadCount)

	case primitive.SortByMostLikes:
		return orderByDesc(fieldLikeCount)

	case primitive.SortByGlobal:
		return fmt.Sprintf("%s, %s, %s", orderByDesc(fieldDownloadCount), orderByDesc(fieldLikeCount),
			orderByDesc(fieldUpdatedAt))

	default:
		return ""
	}
}

// AddLike Add a like to a model by its ID.
func (adapter *modelAdapter) AddLike(model domain.Model) error {
	id := model.Id.Integer()
	version := model.Version

	return adapter.IncrementLikeCount(id, version)
}

// DeleteLike Delete a like to a model by its ID.
func (adapter *modelAdapter) DeleteLike(model domain.Model) error {
	id := model.Id.Integer()
	version := model.Version

	return adapter.DescendLikeCount(id, version)
}
