/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides application services for resource permission management.
package app

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// CanDeleteOrNotFound checks if the user has permission to modify the specified resource or
// the resource is not found for user.
func CanDeleteOrNotFound(
	user primitive.Account, resource domain.Resource, ps ResourcePermissionAppService,
) (bool, error) {
	err := ps.CanDelete(user, resource)
	if err == nil {
		// has permission
		return false, nil
	}

	if !allerror.IsNoPermission(err) {
		return false, err
	}

	return isNoPermissionOrNotFound(user, resource, ps)
}

// CanUpdateOrNotFound checks if the user has permission to update the specified resource or
// the resource is not found for user.
func CanUpdateOrNotFound(
	user primitive.Account, resource domain.Resource, ps ResourcePermissionAppService,
) (bool, error) {
	err := ps.CanUpdate(user, resource)
	if err == nil {
		// has permission
		return false, nil
	}

	if !allerror.IsNoPermission(err) {
		return false, err
	}

	return isNoPermissionOrNotFound(user, resource, ps)
}

func isNoPermissionOrNotFound(
	user primitive.Account, resource domain.Resource, ps ResourcePermissionAppService,
) (bool, error) {
	err := ps.CanRead(user, resource)
	if err == nil {
		// It is no permission when it can read
		return false, allerror.NewNoPermission("no permission", fmt.Errorf("no permission to read"))
	}

	if !allerror.IsNoPermission(err) {
		return false, err
	}

	// It is not found when it can't read,
	return true, nil
}
