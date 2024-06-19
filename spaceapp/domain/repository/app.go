/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repository provides functionality for managing space app repositories.
package repository

import (
	"context"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/spaceapp/domain"
)

// Repository is an interface that defines methods for managing space app repositories.
type Repository interface {
	Add(*domain.SpaceApp) error
	Remove(primitive.Identity) error
	Find(context.Context, *domain.SpaceAppIndex) (domain.SpaceApp, error)
	Save(*domain.SpaceApp) error
	SaveWithBuildLog(*domain.SpaceApp, *domain.SpaceAppBuildLog) error
	FindBySpaceId(context.Context, primitive.Identity) (domain.SpaceApp, error)
	DeleteBySpaceId(primitive.Identity) error
}

// SpaceAppBuildLogAdapter is an interface that defines methods for managing space app build logs.
type SpaceAppBuildLogAdapter interface {
	Find(context.Context, primitive.Identity) (domain.SpaceAppBuildLog, error)
	Save(*domain.SpaceAppBuildLog) error
}
