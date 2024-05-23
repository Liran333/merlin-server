/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package activityrepositoryadapter provides an adapter for the model repository
package activityrepositoryadapter

import (
	"errors"

	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/activity/domain"
	"github.com/openmerlin/merlin-server/activity/domain/repository"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonerror "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

type activityAdapter struct {
	daoImpl
}

func (adapter *activityAdapter) DeleteAll(activity *domain.Activity) error {
	db := adapter.daoImpl.db() // Get the gorm.DB instance for the specific table.
	if db == nil {
		return errors.New("database instance is not initialized")
	}

	resourceIndex := activity.Resource.Index.Integer()
	resourceType := string(activity.Resource.Type)

	if err := db.Where(fieldResourceIndex+" = ? AND "+fieldResourceType+" = ?", resourceIndex, resourceType).Delete(&activityDO{}).Error; err != nil {
		return err
	}

	return nil
}

// Save activities to the database.
func (adapter *activityAdapter) Save(activity *domain.Activity) error {
	db := adapter.daoImpl.db() // Get the gorm.DB instance for the specific table.
	if db == nil {
		return errors.New("database instance is not initialized")
	}

	// Convert domain.Activity to activityDO before saving.
	actDO := activityDO{
		Owner:         activity.Owner.Account(),
		Type:          string(activity.Type),
		Time:          utils.Now(),
		ResourceIndex: activity.Resource.Index.Integer(),
		ResourceType:  string(activity.Resource.Type),
	}

	// Check current count of records for the owner.
	var count int64
	db.Model(&activityDO{}).
		Where(fieldTypeOwner+" = ? AND "+fieldType+" != ?", actDO.Owner, fieldLike).
		Count(&count)

	// If the count is 100 or more, delete the oldest record.
	if count >= config.MaxRecordPerPerson {
		var oldestRecord activityDO
		db.Where(fieldTypeOwner+" = ? AND "+fieldType+" != ?", actDO.Owner, fieldLike).
			Order("time ASC").
			First(&oldestRecord)
		db.Delete(&oldestRecord)
	}

	// Perform the save operation with the converted object.
	result := db.Create(&actDO)
	if result.Error != nil {
		return result.Error // Return the error if the save operation fails.
	}

	return nil
}

// Delete function now using the newly created deleteByOwnerAndIndex function.
func (adapter *activityAdapter) Delete(cmd *domain.Activity) error {
	err := adapter.deleteLikeByOwnerAndIndex(cmd.Owner, cmd.Resource.Index)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return commonerror.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

// Check if the database has a record with specified Type, Owner, and ResourceId and return true if found.
func (adapter *activityAdapter) HasLike(acc primitive.Account, id primitive.Identity) (bool, error) {
	db := adapter.daoImpl.db() // Get the gorm.DB instance.
	if db == nil {
		return false, errors.New("database instance is not initialized")
	}

	var count int64
	// Check if there's any record matching the conditions.
	err := db.Model(&activityDO{}).
		Where("type = ? AND owner = ? AND resource_id = ?", fieldLike, acc, id).
		Count(&count).Error

	// If err is not nil, return the error.
	if err != nil {
		return false, err
	}

	// Return true if a matching record is found, otherwise false.
	return count > 0, nil
}

func order(t primitive.SortType) string {
	if t == nil {
		return ""
	}
	switch t.SortType() {
	case primitive.SortByRecentlyUpdated:
		return orderByDesc(fieldTime)
	default:
		return ""
	}
}

// List retrieves a list of activities based on the provided options.
func (adapter *activityAdapter) List(names []primitive.Account, opt *repository.ListOption) (
	[]domain.Activity, int, error) {
	query := adapter.toQuery(names, opt)

	// Pagination
	if b, offset := opt.Pagination(); b {
		if offset > 0 {
			query = query.Limit(opt.CountPerPage).Offset(offset)
		} else {
			query = query.Limit(opt.CountPerPage)
		}
	}

	// Sorting
	if v := order(opt.SortType); v != "" {
		query = query.Order(v)
	} else {
		//default order
		query = query.Order(orderByDesc(fieldTime))
	}

	var dos []activityDO
	err := query.Find(&dos).Error
	if err != nil {
		return nil, 0, err // Return the error to the caller
	}

	activities := make([]domain.Activity, len(dos))
	for i, do := range dos {
		activity, err := convertToActivityDomain(do)
		if err != nil {
			return nil, 0, err
		}
		activities[i] = activity
	}

	return activities, len(activities), nil
}

// toQuery constructs a GORM DB query with filters based on ListOption.
func (adapter *activityAdapter) toQuery(names []primitive.Account, opt *repository.ListOption) *gorm.DB {
	db := adapter.db() // Assuming this gets a *gorm.DB instance correctly initialized.

	// Accumulate all condition types in a slice
	var conditionTypes []string

	// Add space activity types to conditions if requested
	if opt.Space == primitive.TrueCondition {
		conditionTypes = append(conditionTypes, fieldSpace)
	}

	// Add model activity types to conditions if requested
	if opt.Model == primitive.TrueCondition {
		conditionTypes = append(conditionTypes, fieldModel)
	}

	// Add dataset activity types to conditions if requested
	if opt.Dataset == primitive.TrueCondition {
		conditionTypes = append(conditionTypes, fieldDataset)
	}

	// Begin constructing the query
	query := db.Where("1 = 1")

	// If names are provided, filter by owner
	if len(names) > 0 {
		query = query.Where(fieldTypeOwner+" IN ?", names)
	}

	// If there are any condition types to filter by, add them to the query
	if len(conditionTypes) > 0 {
		query = query.Where(fieldResourceType+" IN ?", conditionTypes)
	}

	if opt.Like == primitive.TrueCondition {
		query = query.Where(fieldType+" IN ?", []string{fieldLike})
	}

	return query
}
