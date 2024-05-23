/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package datasetrepositoryadapter provides an adapter for the dataset repository
package datasetrepositoryadapter

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/datasets/domain"
	"github.com/openmerlin/merlin-server/datasets/domain/repository"
	orgrepo "github.com/openmerlin/merlin-server/organization/domain/repository"
)

const (
	// README represents the constant value for the README file.
	README = "README"
	readme = "readme"
)

type datasetAdapter struct {
	daoImpl
}

// Add adds a new dataset to the database.
func (adapter *datasetAdapter) Add(dataset *domain.Dataset) error {
	do := toDatasetDO(dataset)

	v := adapter.db().Create(&do)

	return v.Error
}

// FindByName finds a dataset by its name.
func (adapter *datasetAdapter) FindByName(index *domain.DatasetIndex) (domain.Dataset, error) {
	do := datasetDO{Owner: index.Owner.Account(), Name: index.Name.MSDName()}

	if err := adapter.GetLowerDatasetName(&do, &do); err != nil {
		return domain.Dataset{}, err
	}

	return do.toDataset(), nil
}

// FindById finds a dataset by its ID.
func (adapter *datasetAdapter) FindById(datasetId primitive.Identity) (domain.Dataset, error) {
	do := datasetDO{Id: datasetId.Integer()}

	if err := adapter.GetByPrimaryKey(&do); err != nil {
		return domain.Dataset{}, err
	}

	return do.toDataset(), nil
}

// Delete deletes a dataset by its ID.
func (adapter *datasetAdapter) Delete(datasetId primitive.Identity) error {
	return adapter.DeleteByPrimaryKey(
		&datasetDO{Id: datasetId.Integer()},
	)
}

// Save updates an existing dataset in the database.
func (adapter *datasetAdapter) Save(dataset *domain.Dataset) error {
	do := toDatasetDO(dataset)
	do.Version += 1

	v := adapter.db().Model(
		&datasetDO{Id: dataset.Id.Integer()},
	).Where(
		equalQuery(fieldVersion), dataset.Version,
	).Select(`*`).Omit(fieldTask, fieldSize, fieldLanguage, fieldDomain).Updates(&do)

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

// List retrieves a list of dataset based on the provided options.
func (adapter *datasetAdapter) List(opt *repository.ListOption, login primitive.Account, member orgrepo.OrgMember) ([]repository.DatasetSummary, int, error) {
	query := adapter.toQuery(opt)

	if opt.Visibility != nil {
		if login != nil {
			members, err := member.GetByUser(login.Account())
			if err != nil {
				return nil, 0, err
			}
			orgNames := make([]string, 0, len(members))
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

	var dos []datasetDO

	err := query.Find(&dos).Error
	if err != nil || len(dos) == 0 {
		return nil, 0, nil
	}

	r := make([]repository.DatasetSummary, len(dos))
	for i := range dos {
		r[i] = dos[i].toDatasetSummary()
	}

	return r, int(total), nil
}

// Count counts the number of dataset based on the provided options.
func (adapter *datasetAdapter) Count(opt *repository.ListOption) (int, error) {
	var total int64
	err := adapter.toQuery(opt).Count(&total).Error

	return int(total), err
}

func (adapter *datasetAdapter) toQuery(opt *repository.ListOption) *gorm.DB {
	db := adapter.db()

	if opt.Name != "" {
		_, arg := likeFilter(fieldName, opt.Name)

		if !opt.ExcludeFullname {
			query2, arg2 := likeFilter(fieldFullName, opt.Name)
			db = db.Where(db.Where(gorm.Expr("CONCAT("+fieldOwner+", '/', "+fieldName+") ilike ?", arg)).Or(query2, arg2))
		} else {
			db = db.Where(db.Where(gorm.Expr("CONCAT("+fieldOwner+", '/', "+fieldName+") ilike ?", arg)))
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

	if v := opt.Labels.Task; v != nil && v.Len() > 0 {
		query, arg := intersectionFilter(fieldTask, v.UnsortedList())

		db = db.Where(query, arg)
	}

	if v := opt.Labels.Size; v != "" {
		db = db.Where(equalQuery(fieldSize), v)
	}

	if v := opt.Labels.Language; v != nil && v.Len() > 0 {
		query, arg := intersectionFilter(fieldLanguage, v.UnsortedList())

		db = db.Where(query, arg)
	}

	if v := opt.Labels.Domain; v != nil && v.Len() > 0 {
		query, arg := intersectionFilter(fieldDomain, v.UnsortedList())

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

// AddLike Add a like to a dataset by its ID.
func (adapter *datasetAdapter) AddLike(dataset domain.Dataset) error {
	id := dataset.Id.Integer()
	version := dataset.Version

	return adapter.IncrementLikeCount(id, version)
}

// DeleteLike Delete a like to a dataset by its ID.
func (adapter *datasetAdapter) DeleteLike(dataset domain.Dataset) error {
	id := dataset.Id.Integer()
	version := dataset.Version

	return adapter.DescendLikeCount(id, version)
}
