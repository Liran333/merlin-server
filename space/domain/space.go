/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package domain

import (
	"k8s.io/apimachinery/pkg/util/sets"

	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	spaceprimitive "github.com/openmerlin/merlin-server/space/domain/primitive"
	"github.com/openmerlin/merlin-server/space/domain/securestorage"
)

const (
	variablePath = "variable/"
	secretePath  = "secret/"
)

// Space represents a space with its associated properties and methods.
type Space struct {
	coderepo.CodeRepo

	SDK      spaceprimitive.SDK
	Desc     primitive.MSDDesc
	Labels   SpaceLabels
	Fullname primitive.MSDFullname
	Hardware spaceprimitive.Hardware

	LocalCmd      string
	LocalEnvInfo  string
	Version       int
	CreatedAt     int64
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
}

// ResourceType returns the type of the model resource.
func (m *Space) ResourceType() primitive.ObjType {
	return primitive.ObjTypeSpace
}

// SpaceLabels represents labels associated with a space.
type SpaceLabels struct {
	Task       string           // task label
	Others     sets.Set[string] // other labels
	Frameworks sets.Set[string] // framework labels
}

// SpaceIndex represents an index for spaces in the code repository.
type SpaceIndex = coderepo.CodeRepoIndex

type SpaceVariable struct {
	Id      primitive.Identity
	SpaceId primitive.Identity
	Name    primitive.MSDName
	Desc    primitive.MSDDesc
	Value   spaceprimitive.ENVValue

	CreatedAt int64
	UpdatedAt int64
}

// NewSpaceVariableVault return a space env secret vault by space variable
func NewSpaceVariableVault(variable *SpaceVariable) securestorage.SpaceEnvSecret {
	return securestorage.SpaceEnvSecret{
		Path:  variablePath + variable.SpaceId.Identity(),
		Name:  variable.Name.MSDName(),
		Value: variable.Value.ENVValue(),
	}
}

// GetVariablePath return vault space variable path
func (variable *SpaceVariable) GetVariablePath() string {
	return variablePath + variable.SpaceId.Identity()
}

type SpaceSecret struct {
	Id      primitive.Identity
	SpaceId primitive.Identity
	Name    primitive.MSDName
	Desc    primitive.MSDDesc
	Value   spaceprimitive.ENVValue

	CreatedAt int64
	UpdatedAt int64
}

// NewSpaceSecretVault return a space env secret vault by space secret
func NewSpaceSecretVault(secret *SpaceSecret) securestorage.SpaceEnvSecret {
	return securestorage.SpaceEnvSecret{
		Path:  secretePath + secret.SpaceId.Identity(),
		Name:  secret.Name.MSDName(),
		Value: secret.Value.ENVValue(),
	}
}

// GetSecretPath return vault space secret path
func (secret *SpaceSecret) GetSecretPath() string {
	return secretePath + secret.SpaceId.Identity()
}
