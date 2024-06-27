/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"io"
	"math/rand"

	sdk "github.com/openmerlin/merlin-sdk/space"
	"k8s.io/apimachinery/pkg/util/sets"

	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space/domain"
	spaceprimitive "github.com/openmerlin/merlin-server/space/domain/primitive"
	"github.com/openmerlin/merlin-server/space/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

// CmdToCreateSpace is a struct used to create a space.
type CmdToCreateSpace struct {
	coderepoapp.CmdToCreateRepo

	SDK          spaceprimitive.SDK
	Desc         primitive.MSDDesc
	Fullname     primitive.MSDFullname
	Hardware     spaceprimitive.Hardware
	HardwareType string
	BaseImage    spaceprimitive.BaseImage
	AvatarId     primitive.Avatar
}

func (cmd *CmdToCreateSpace) toSpace() domain.Space {
	if cmd.AvatarId != nil && cmd.AvatarId.URL() == "" {
		configAvatarId := config.avatarIdsSet.UnsortedList()[rand.Intn(len(config.avatarIdsSet))] // #nosec G404
		cmd.AvatarId = primitive.CreateAvatar(configAvatarId)
	}

	label := domain.SpaceLabels{
		Framework: cmd.BaseImage.Type(),
		Licenses:  cmd.License,
	}

	s := domain.Space{
		SDK:       cmd.SDK,
		Desc:      cmd.Desc,
		Hardware:  cmd.Hardware,
		BaseImage: cmd.BaseImage,
		Fullname:  cmd.Fullname,
		AvatarId:  cmd.AvatarId,
		Labels:    label,
	}

	return s
}

// CmdToUpdateSpace is a struct used to update a space.
type CmdToUpdateSpace struct {
	coderepoapp.CmdToUpdateRepo

	SDK      spaceprimitive.SDK
	Desc     primitive.MSDDesc
	Fullname primitive.MSDFullname
	Hardware spaceprimitive.Hardware
	AvatarId primitive.Avatar
}

func (cmd *CmdToUpdateSpace) toSpace(space *domain.Space) (b bool) {
	if v := cmd.SDK; v != nil && v != space.SDK {
		space.SDK = v
		b = true
	}

	if v := cmd.Desc; v != nil && v != space.Desc {
		space.Desc = v
		b = true
	}

	if v := cmd.Fullname; v != nil && v != space.Fullname {
		space.Fullname = v
		b = true
	}

	if v := cmd.Hardware; v != nil && v != space.Hardware {
		space.Hardware = v
		space.Labels.HardwareType = space.GetHardwareType()
		b = true
	}

	if v := cmd.AvatarId; v != nil && v.URL() != space.AvatarId.URL() {
		space.AvatarId = v
		b = true
	}

	if b {
		space.UpdatedAt = utils.Now()
	}

	return
}

// CmdToDisableSpace is a struct used to disable a space.
type CmdToDisableSpace struct {
	Disable       bool
	DisableReason primitive.DisableReason
}

func (cmd *CmdToDisableSpace) toSpace(space *domain.Space) {
	space.Disable = cmd.Disable
	space.DisableReason = cmd.DisableReason
	space.UpdatedAt = utils.Now()
}

// SpaceDTO is a struct used to represent a space data transfer object.
type SpaceDTO struct {
	Id            string         `json:"id"`
	SDK           string         `json:"sdk"`
	Name          string         `json:"name"`
	Desc          string         `json:"desc"`
	Owner         string         `json:"owner"`
	Labels        SpaceLabelsDTO `json:"labels"`
	Fullname      string         `json:"fullname"`
	AvatarId      string         `json:"space_avatar_id"`
	Hardware      string         `json:"hardware"`
	BaseImage     string         `json:"base_image"`
	CreatedAt     int64          `json:"created_at"`
	UpdatedAt     int64          `json:"updated_at"`
	LikeCount     int            `json:"like_count"`
	LocalCMD      string         `json:"local_cmd"`
	LocalEnvInfo  string         `json:"local_env_info"`
	Visibility    string         `json:"visibility"`
	DownloadCount int            `json:"download_count"`
	VisitCount    int            `json:"visit_count"`
	Disable       bool           `json:"disable"`
	DisableReason string         `json:"disable_reason"`
	Exception     string         `json:"exception"`

	IsNpu              bool `json:"is_npu"`
	CompPowerAllocated bool `json:"comp_power_allocated"`
	NoApplicationFile  bool `json:"no_application_file"`
}

// SpaceLabelsDTO is a struct used to represent labels of a space.
type SpaceLabelsDTO struct {
	Task      string   `json:"task"`
	License   []string `json:"license"`
	Framework string   `json:"framework"`
}

func toSpaceLabelsDTO(space *domain.Space) SpaceLabelsDTO {
	labels := &space.Labels

	return SpaceLabelsDTO{
		Task:      labels.Task.Task(),
		License:   space.License.License(),
		Framework: labels.Framework,
	}
}

func toSpaceDTO(space *domain.Space) SpaceDTO {
	dto := SpaceDTO{
		Id:            space.Id.Identity(),
		SDK:           space.SDK.SDK(),
		Name:          space.Name.MSDName(),
		Owner:         space.Owner.Account(),
		Labels:        toSpaceLabelsDTO(space),
		Hardware:      space.Hardware.Hardware(),
		BaseImage:     space.BaseImage.BaseImage(),
		CreatedAt:     space.CreatedAt,
		UpdatedAt:     space.UpdatedAt,
		LikeCount:     space.LikeCount,
		AvatarId:      space.AvatarId.URL(),
		Visibility:    space.Visibility.Visibility(),
		DownloadCount: space.DownloadCount,
		VisitCount:    space.VisitCount,
		LocalCMD:      space.GetLocalCmd(),
		LocalEnvInfo:  space.GetLocalEnvInfo(),
		Disable:       space.Disable,
		DisableReason: space.DisableReason.DisableReason(),
		Exception:     space.Exception.Exception(),

		IsNpu:              space.Hardware.IsNpu(),
		CompPowerAllocated: space.CompPowerAllocated,
		NoApplicationFile:  space.NoApplicationFile,
	}

	if space.Desc != nil {
		dto.Desc = space.Desc.MSDDesc()
	}

	if space.Fullname != nil {
		dto.Fullname = space.Fullname.MSDFullname()
	}

	return dto
}

func toSpaceSummary(spaceDTO *SpaceDTO) repository.SpaceSummary {
	return repository.SpaceSummary{
		Id:            spaceDTO.Id,
		Name:          spaceDTO.Name,
		Desc:          spaceDTO.Desc,
		Owner:         spaceDTO.Owner,
		Fullname:      spaceDTO.Fullname,
		BaseImage:     spaceDTO.BaseImage,
		AvatarId:      spaceDTO.AvatarId,
		UpdatedAt:     spaceDTO.UpdatedAt,
		LikeCount:     spaceDTO.LikeCount,
		DownloadCount: spaceDTO.DownloadCount,
		VisitCount:    spaceDTO.VisitCount,
		Disable:       spaceDTO.Disable,
		DisableReason: spaceDTO.DisableReason,
		Labels: domain.SpaceLabels{
			Task:      spaceprimitive.CreateTask(spaceDTO.Labels.Task),
			Licenses:  primitive.CreateLicense(spaceDTO.Labels.License),
			Framework: spaceDTO.Labels.Framework,
		},
		IsNpu:              spaceDTO.IsNpu,
		Exception:          spaceDTO.Exception,
		CompPowerAllocated: spaceDTO.CompPowerAllocated,
		NoApplicationFile:  spaceDTO.NoApplicationFile,
	}
}

// SpacesDTO represents the data transfer object for spaces.
type SpacesDTO struct {
	Total  int                       `json:"total"`
	Spaces []repository.SpaceSummary `json:"spaces"`
}

// CmdToListSpaces is a command to list spaces with repository.ListOption options.
type CmdToListSpaces = repository.ListOption

func toSpaceMetaDTO(space *domain.Space) sdk.SpaceMetaDTO {
	dto := sdk.SpaceMetaDTO{
		Id:         space.Id.Identity(),
		SDK:        space.SDK.SDK(),
		Name:       space.Name.MSDName(),
		Owner:      space.Owner.Account(),
		Hardware:   space.Hardware.Hardware(),
		BaseImage:  space.BaseImage.BaseImage(),
		Visibility: space.CodeRepo.Visibility.Visibility(),
		Disable:    space.Disable,
	}
	return dto
}

// SpaceModelDTO
type SpaceModelDTO struct {
	Owner         string `json:"owner"`
	Name          string `json:"name"`
	AvatarId      string `json:"avatar_id"`
	UpdatedAt     int64  `json:"updated_at"`
	LikeCount     int    `json:"like_count"`
	DownloadCount int    `json:"download_count"`
	VisitCount    int    `json:"visit_count"`
}

// SpaceIdModelDTO
type SpaceIdModelDTO struct {
	SpaceId []string `json:"space_id"`
}

// CmdToCreateSpaceVariable is a struct used to create a space variable.
type CmdToCreateSpaceVariable struct {
	Name  spaceprimitive.ENVName
	Desc  primitive.MSDDesc
	Value spaceprimitive.ENVValue
}

// CmdToUpdateSpaceVariable is a struct used to update a space variable.
type CmdToUpdateSpaceVariable struct {
	Desc  primitive.MSDDesc
	Value spaceprimitive.ENVValue
}

func (cmd *CmdToUpdateSpaceVariable) toSpaceVariable(spaceVariable *domain.SpaceVariable) (b bool) {
	if v := cmd.Desc; v != nil && v != spaceVariable.Desc {
		spaceVariable.Desc = v
		b = true
	}

	if v := cmd.Value; v != nil && v != spaceVariable.Value {
		spaceVariable.Value = v
		b = true
	}

	return
}

// SpaceVariableSecretDTO represents the data transfer object for spaces variable and secret.
type SpaceVariableSecretDTO struct {
	SpaceVariableSecret []repository.SpaceVariableSecretSummary `json:"space_variable_secret"`
}

// CmdToCreateSpaceSecret is a struct used to create a space secret.
type CmdToCreateSpaceSecret struct {
	Name  spaceprimitive.ENVName
	Desc  primitive.MSDDesc
	Value spaceprimitive.ENVValue
}

// CmdToUpdateSpaceSecret is a struct used to update a space secret.
type CmdToUpdateSpaceSecret struct {
	Desc  primitive.MSDDesc
	Value spaceprimitive.ENVValue
}

func (cmd *CmdToUpdateSpaceSecret) toSpaceSecret(spaceSecret *domain.SpaceSecret) (b bool) {
	if v := cmd.Desc; v != nil && v != spaceSecret.Desc {
		spaceSecret.Desc = v
		b = true
	}

	if v := cmd.Value; v != nil && v != spaceSecret.Value {
		spaceSecret.Value = v
		b = true
	}

	return
}

// CmdToUpdateStatistics is to update download count
type CmdToUpdateStatistics struct {
	DownloadCount int `json:"download_count"`
	VisitCount    int `json:"visit_count"`
}

// CmdToResetLabels is a type alias for domain.SpaceLabels, representing a command to reset space labels.
type CmdToResetLabels struct {
	Task     spaceprimitive.Task
	Licenses sets.Set[string] // license label
}

// CmdToNotifyUpdateCode is to update no application file and commitId
type CmdToNotifyUpdateCode struct {
	CommitId string
	HasHtml  bool
	HasApp   bool
}

// CmdToUploadCover is to update no application file and commitId
type CmdToUploadCover struct {
	Image    io.Reader
	User     primitive.Account
	FileName string
}

type SpaceCoverDTO struct {
	URL string `json:"url"`
}

func toSpaceCoverDTO(u string) SpaceCoverDTO {
	return SpaceCoverDTO{
		URL: u,
	}
}
