/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repositoryadapter provides an adapter implementation for working with the repository of space applications.
package repositoryadapter

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/domain/repository"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/spaceapp/domain"
)

type dao interface {
	DB() *gorm.DB
	EqualQuery(field string) string
	NotEqualQuery(field string) string
	IsRecordExists(err error) bool
	GetRecord(ctx context.Context, filter, result interface{}) error
	GetByPrimaryKey(ctx context.Context, row interface{}) error
	DeleteByPrimaryKey(ctx context.Context, row interface{}) error
}

type appRepositoryAdapter struct {
	dao dao
}

// Add adds a space application to the repository.
func (adapter *appRepositoryAdapter) Add(m *domain.SpaceApp) error {
	do := toSpaceAppDO(m)

	err := adapter.dao.DB().Create(&do).Error

	if err != nil && adapter.dao.IsRecordExists(err) {
		return repository.NewErrorDuplicateCreating(
			errors.New("space app exists"),
		)
	}
	return err
}

func (adapter *appRepositoryAdapter) Remove(spaceId primitive.Identity) error {
	return adapter.dao.DB().Where(
		adapter.dao.EqualQuery(fieldSpaceId), spaceId.Identity(),
	).Delete(
		spaceappDO{},
	).Error
}

// FindBySpaceId finds a space application in the repository based on the space ID.
func (adapter *appRepositoryAdapter) FindBySpaceId(
	ctx context.Context, spaceId primitive.Identity) (domain.SpaceApp, error) {
	do := spaceappDO{SpaceId: spaceId.Integer()}

	// It must new a new DO, otherwise the sql statement will include duplicate conditions.
	result := spaceappDO{}

	if err := adapter.dao.GetRecord(ctx, &do, &result); err != nil {
		return domain.SpaceApp{}, err
	}

	return result.toSpaceApp(), nil
}

// Find finds a space application in the repository based on the space app index.
func (adapter *appRepositoryAdapter) Find(ctx context.Context, index *domain.SpaceAppIndex) (domain.SpaceApp, error) {
	do := spaceappDO{SpaceId: index.SpaceId.Integer(), CommitId: index.CommitId}

	// It must new a new DO, otherwise the sql statement will include duplicate conditions.
	result := spaceappDO{}

	if err := adapter.dao.GetRecord(ctx, &do, &result); err != nil {
		return domain.SpaceApp{}, err
	}

	return result.toSpaceApp(), nil
}

// Save saves a space application in the repository.
func (adapter *appRepositoryAdapter) Save(m *domain.SpaceApp) error {
	do := toSpaceAppDO(m)
	do.Version += 1

	v := adapter.dao.DB().Model(
		&spaceappDO{Id: m.Id.Integer()},
	).Where(
		adapter.dao.EqualQuery(fieldVersion), m.Version,
	).Select(`*`).Omit(fieldAllBuildLog).Updates(&do)

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

// SaveWithBuildLog saves a space application and build log in the repository.
func (adapter *appRepositoryAdapter) SaveWithBuildLog(m *domain.SpaceApp, log *domain.SpaceAppBuildLog) error {
	do := toSpaceAppDO(m)
	do.Version += 1
	do.AllBuildLog = log.Logs

	v := adapter.dao.DB().Model(
		&spaceappDO{Id: m.Id.Integer()},
	).Where(
		adapter.dao.EqualQuery(fieldVersion), m.Version,
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

// DeleteBySpaceId delete space app by the space ID.
func (adapter *appRepositoryAdapter) DeleteBySpaceId(spaceId primitive.Identity) error {
	return adapter.dao.DB().Where(equalQuery(fieldSpaceId), spaceId.Identity()).Delete(&spaceappDO{}).Error
}
