/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package domain

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"

	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	spaceprimitive "github.com/openmerlin/merlin-server/space/domain/primitive"
	"github.com/openmerlin/merlin-server/space/domain/securestorage"
)

const (
	variablePath = "variable/"
	secretePath  = "secret/"

	computilityTypeNpu = "npu"
	computilityTypeCpu = "cpu"
)

// Space represents a space with its associated properties and methods.
type Space struct {
	coderepo.CodeRepo

	SDK           spaceprimitive.SDK
	Desc          primitive.MSDDesc
	Labels        SpaceLabels
	Fullname      primitive.MSDFullname
	Hardware      spaceprimitive.Hardware
	AvatarId      primitive.Avatar
	BaseImage     spaceprimitive.BaseImage
	LocalCmd      string
	LocalEnvInfo  string
	Version       int
	CreatedAt     int64
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
	VisitCount    int

	Disable       bool
	DisableReason primitive.DisableReason
	Exception     primitive.Exception

	CompPowerAllocated   bool
	NoApplicationFile    bool
	CommitId             string
	IsDiscussionDisabled bool
}

// ResourceType returns the type of the model resource.
func (m *Space) ResourceType() primitive.ObjType {
	return primitive.ObjTypeSpace
}

// IsDisable checks if the space is disable.
func (m *Space) IsDisable() bool {
	return m.Disable
}

func (m *Space) DiscussionDisabled() bool {
	return m.IsDiscussionDisabled
}

func (m *Space) CloseDiscussion() error {
	if m.IsDiscussionDisabled {
		return allerror.New(
			allerror.ErrorCodeDiscussionDisabled,
			"failed to close discussion",
			xerrors.Errorf("discussion is closed"),
		)
	}

	m.IsDiscussionDisabled = true

	return nil
}

func (m *Space) ReopenDiscussion() error {
	if !m.IsDiscussionDisabled {
		return allerror.New(
			allerror.ErrorCodeDiscussionEnabled,
			"failed to reopen discussion",
			xerrors.Errorf("discussion is open"),
		)
	}

	m.IsDiscussionDisabled = false

	return nil
}

// IsNoApplicationFile checks if the space is valid.
func (m *Space) IsNoApplicationFile() bool {
	return m.NoApplicationFile
}

// GetLocalCmd return if the space is public.
func (m *Space) GetLocalCmd() string {
	if m.IsPublic() {
		return m.LocalCmd
	}
	return ""
}

// GetLocalEnvInfo return if the space is public.
func (m *Space) GetLocalEnvInfo() string {
	if m.IsPublic() {
		return m.LocalEnvInfo
	}
	return ""
}

// SetSpaceCommitId for update space commitId.
func (m *Space) SetSpaceCommitId(commitId string) {
	m.CommitId = commitId
}

// SetNoApplicationFile for set NoApplicationFile and Exception.
func (m *Space) SetNoApplicationFile(hasHtml, hasApp bool) {
	m.NoApplicationFile = true
	if (m.SDK == spaceprimitive.StaticSdk) && hasHtml {
		m.NoApplicationFile = false
	}
	if (m.SDK == spaceprimitive.GradioSdk) && hasApp {
		m.NoApplicationFile = false
	}
	if m.NoApplicationFile {
		m.Exception = primitive.CreateException(primitive.NoApplicationFile)
		return
	}
	if !m.NoApplicationFile && m.Exception == primitive.ExceptionNoApplicationFile {
		m.Exception = primitive.CreateException("")
	}
}

// SpaceLabels represents labels associated with a space.
type SpaceLabels struct {
	Task         spaceprimitive.Task // task label
	Licenses     primitive.License   // license label
	HardwareType string              // hardware_type
	Framework    string              // framework
}

// SpaceIndex represents an index for spaces in the code repository.
type SpaceIndex = coderepo.CodeRepoIndex

// SpaceVariable represents a variable associated with a space.
type SpaceVariable struct {
	Id      primitive.Identity
	SpaceId primitive.Identity
	Name    spaceprimitive.ENVName
	Desc    primitive.MSDDesc
	Value   spaceprimitive.ENVValue

	CreatedAt int64
	UpdatedAt int64
}

// NewSpaceVariableVault return a space env secret vault by space variable
func NewSpaceVariableVault(variable *SpaceVariable) securestorage.SpaceEnvSecret {
	return securestorage.SpaceEnvSecret{
		Path:  variablePath + variable.SpaceId.Identity(),
		Name:  variable.Name.ENVName(),
		Value: variable.Value.ENVValue(),
	}
}

// GetVariablePath return vault space variable path
func (variable *SpaceVariable) GetVariablePath() string {
	return variablePath + variable.SpaceId.Identity()
}

// SpaceSecret represents a secret associated with a space.
type SpaceSecret struct {
	Id      primitive.Identity
	SpaceId primitive.Identity
	Name    spaceprimitive.ENVName
	Desc    primitive.MSDDesc
	Value   spaceprimitive.ENVValue

	CreatedAt int64
	UpdatedAt int64
}

// NewSpaceSecretVault return a space env secret vault by space secret
func NewSpaceSecretVault(secret *SpaceSecret) securestorage.SpaceEnvSecret {
	return securestorage.SpaceEnvSecret{
		Path:  secretePath + secret.SpaceId.Identity(),
		Name:  secret.Name.ENVName(),
		Value: secret.Value.ENVValue(),
	}
}

// GetSecretPath return vault space secret path
func (secret *SpaceSecret) GetSecretPath() string {
	return secretePath + secret.SpaceId.Identity()
}

// GetComputeType returns the compute type of the Space.
func (s *Space) GetComputeType() primitive.ComputilityType {
	if s.Hardware.IsNpu() {
		return primitive.CreateComputilityType(computilityTypeNpu)
	} else if s.Hardware.IsCpu() {
		return primitive.CreateComputilityType(computilityTypeCpu)
	}

	return nil
}

// GetQuotaCount returns the quota count of the Space.
func (s *Space) GetQuotaCount() int {
	if s.Hardware.IsNpu() {
		return 1
	} else if s.Hardware.IsCpu() {
		return 0
	}

	return 0
}

func (s *Space) GetHardwareType() string {
	if s.Hardware.IsNpu() {
		return "npu"
	} else if s.Hardware.IsCpu() {
		return "cpu"
	}
	return ""
}

// IsValidHardwareType checks if the provided hardware type string is a valid hardware type.
func IsValidHardwareType(h string) bool {
	return strings.ToLower(h) == computilityTypeNpu || strings.ToLower(h) == computilityTypeCpu
}

type CoverInfo struct {
	User        primitive.Account
	FileName    string
	Path        string
	Bucket      string
	CdnEndpoint string
}

func (c *CoverInfo) GetCoverURL() string {
	return fmt.Sprintf("%s%s/%s/%s", c.CdnEndpoint, c.Path, c.User.Account(), c.FileName)
}

func (c *CoverInfo) GetObsPath() string {
	return fmt.Sprintf("%s/%s/%s", c.Path, c.User.Account(), c.FileName)
}
