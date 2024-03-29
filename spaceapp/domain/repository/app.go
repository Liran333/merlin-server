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
	Find(*domain.SpaceAppIndex) (domain.SpaceApp, error)
	Save(*domain.SpaceApp) error
	FindBySpaceId(primitive.Identity) (domain.SpaceApp, error)
	DeleteBySpaceId(primitive.Identity) error
}

type SpaceAppBuildLogAdapter interface {
	Save(*domain.SpaceAppBuildLog) error
}
