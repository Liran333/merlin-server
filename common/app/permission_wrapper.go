/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides application services for resource permission management.
package app

import (
	"context"
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// CanDeleteOrNotFound checks if the user has permission to modify the specified resource or
// the resource is not found for user.
func CanDeleteOrNotFound(
	ctx context.Context, user primitive.Account, resource domain.Resource, ps ResourcePermissionAppService,
) (bool, error) {
	err := ps.CanDelete(ctx, user, resource)
	if err == nil {
		// has permission
		return false, nil
	}

	if !allerror.IsNoPermission(err) {
		return false, err
	}

	return isNoPermissionOrNotFound(ctx, user, resource, ps)
}

// CanUpdateOrNotFound checks if the user has permission to update the specified resource or
// the resource is not found for user.
func CanUpdateOrNotFound(
	ctx context.Context, user primitive.Account, resource domain.Resource, ps ResourcePermissionAppService,
) (bool, error) {
	err := ps.CanUpdate(ctx, user, resource)
	if err == nil {
		// has permission
		return false, nil
	}

	if !allerror.IsNoPermission(err) {
		return false, err
	}

	return isNoPermissionOrNotFound(ctx, user, resource, ps)
}

func isNoPermissionOrNotFound(
	ctx context.Context, user primitive.Account, resource domain.Resource, ps ResourcePermissionAppService,
) (bool, error) {
	err := ps.CanRead(ctx, user, resource)
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

// CanReadOrNotFound checks if the user has permission to read the specified resource or
// the resource is not found for user.
func CanReadOrNotFound(
	ctx context.Context, user primitive.Account, resource domain.Resource, ps ResourcePermissionAppService,
) (bool, error) {
	err := ps.CanRead(ctx, user, resource)
	if err == nil {
		// has permission
		return false, nil
	}

	return true, allerror.NewNotFound(allerror.ErrorCodeRepoNotFound, "not found", err)
}
