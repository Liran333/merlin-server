/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repository provides functionality for managing space app repositories.
package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/spaceapp/domain"
)

// Repository is an interface that defines methods for managing space app repositories.
type Repository interface {
	Add(*domain.SpaceApp) error
	Remove(primitive.Identity) error
	Find(*domain.SpaceAppIndex) (domain.SpaceApp, error)
	Save(*domain.SpaceApp) error
	SaveWithBuildLog(*domain.SpaceApp, *domain.SpaceAppBuildLog) error
	FindBySpaceId(primitive.Identity) (domain.SpaceApp, error)
	DeleteBySpaceId(primitive.Identity) error
}

// SpaceAppBuildLogAdapter is an interface that defines methods for managing space app build logs.
type SpaceAppBuildLogAdapter interface {
	Find(primitive.Identity) (domain.SpaceAppBuildLog, error)
	Save(*domain.SpaceAppBuildLog) error
}
