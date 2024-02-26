/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package platform

import "github.com/openmerlin/merlin-server/user/domain"

// User is an interface for user related operations.
type User interface {
	Create(*domain.UserCreateCmd) (domain.User, error)
	Delete(user *domain.User) error
}
