/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"math/rand"

	sdk "github.com/openmerlin/merlin-sdk/space"

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

	SDK      spaceprimitive.SDK
	Desc     primitive.MSDDesc
	Fullname primitive.MSDFullname
	Hardware spaceprimitive.Hardware
	AvatarId primitive.AvatarId
}

func (cmd *CmdToCreateSpace) toSpace() domain.Space {
	if cmd.AvatarId != nil && cmd.AvatarId.AvatarId() == "" {
		cmd.AvatarId = primitive.CreateAvatarId(config.avatarIdsSet.UnsortedList()[rand.Intn(len(config.avatarIdsSet))]) // #nosec G404
	}

	return domain.Space{
		SDK:      cmd.SDK,
		Desc:     cmd.Desc,
		Hardware: cmd.Hardware,
		Fullname: cmd.Fullname,
		AvatarId: cmd.AvatarId,
	}
}

// CmdToUpdateSpace is a struct used to update a space.
type CmdToUpdateSpace struct {
	coderepoapp.CmdToUpdateRepo

	SDK      spaceprimitive.SDK
	Desc     primitive.MSDDesc
	Fullname primitive.MSDFullname
	Hardware spaceprimitive.Hardware
	AvatarId primitive.AvatarId
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
		b = true
	}

	if v := cmd.AvatarId; v != nil && v != space.AvatarId {
		space.AvatarId = v
		b = true
	}

	if b {
		space.UpdatedAt = utils.Now()
	}

	return
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
	CreatedAt     int64          `json:"created_at"`
	UpdatedAt     int64          `json:"updated_at"`
	LikeCount     int            `json:"like_count"`
	LocalCMD      string         `json:"local_cmd"`
	LocalEnvInfo  string         `json:"local_env_info"`
	Visibility    string         `json:"visibility"`
	DownloadCount int            `json:"download_count"`
}

// SpaceLabelsDTO is a struct used to represent labels of a space.
type SpaceLabelsDTO struct {
	Task       string   `json:"task"`
	Others     []string `json:"others"`
	License    string   `json:"license"`
	Frameworks []string `json:"frameworks"`
}

func toSpaceLabelsDTO(space *domain.Space) SpaceLabelsDTO {
	labels := &space.Labels

	return SpaceLabelsDTO{
		Task:       labels.Task,
		Others:     labels.Others.UnsortedList(),
		License:    space.License.License(),
		Frameworks: labels.Frameworks.UnsortedList(),
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
		CreatedAt:     space.CreatedAt,
		UpdatedAt:     space.UpdatedAt,
		LikeCount:     space.LikeCount,
		AvatarId:      space.AvatarId.AvatarId(),
		Visibility:    space.Visibility.Visibility(),
		DownloadCount: space.DownloadCount,
		LocalCMD:      space.LocalCmd,
		LocalEnvInfo:  space.LocalEnvInfo,
	}

	if space.Desc != nil {
		dto.Desc = space.Desc.MSDDesc()
	}

	if space.Fullname != nil {
		dto.Fullname = space.Fullname.MSDFullname()
	}

	return dto
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
		Visibility: space.CodeRepo.Visibility.Visibility(),
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
}

// CmdToCreateSpaceVariable is a struct used to create a space variable.
type CmdToCreateSpaceVariable struct {
	Name  primitive.MSDName
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
	Name  primitive.MSDName
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

// CmdToUpdateStatistics is a type alias for domain.ModelLabels, representing a command to update model statistics.
type CmdToUpdateStatistics struct {
	DownloadCount int `json:"download_count"`
}
