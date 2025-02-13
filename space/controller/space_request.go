/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openmerlin/merlin-sdk/space"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openmerlin/merlin-server/common/controller"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/domain"
	"github.com/openmerlin/merlin-server/space/app"
	spaceprimitive "github.com/openmerlin/merlin-server/space/domain/primitive"
	"github.com/openmerlin/merlin-server/space/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

const (
	firstPage          = 1
	labelSpliter       = ","
	repoNameSplitedLen = 2
)

type reqToCreateSpace struct {
	SDK        string `json:"sdk"        required:"true"`
	Name       string `json:"name"       required:"true"`
	Desc       string `json:"desc"`
	Owner      string `json:"owner"      required:"true"`
	License    string `json:"license"    required:"true"`
	Hardware   string `json:"hardware"   required:"true"`
	BaseImage  string `json:"base_image" required:"true"`
	Fullname   string `json:"fullname"`
	Visibility string `json:"visibility" required:"true"`
	AvatarId   string `json:"space_avatar_id"`
}

func (req *reqToCreateSpace) action() string {
	return fmt.Sprintf("create space of %s/%s", req.Owner, req.Name)
}

func (req *reqToCreateSpace) toCmd() (cmd app.CmdToCreateSpace, err error) {
	if cmd.SDK, err = spaceprimitive.NewSDK(req.SDK); err != nil {
		err = xerrors.Errorf("invalid sdk: %w", err)
		return
	}

	if cmd.Name, err = primitive.NewMSDName(req.Name); err != nil {
		err = xerrors.Errorf("invalid name: %w", err)
		return
	}

	if cmd.Desc, err = primitive.NewMSDDesc(req.Desc); err != nil {
		err = xerrors.Errorf("invalid desc: %w", err)
		return
	}

	if cmd.Owner, err = primitive.NewAccount(req.Owner); err != nil {
		err = xerrors.Errorf("invalid owner: %w", err)
		return
	}

	if cmd.License, err = primitive.NewLicense(req.License); err != nil {
		err = xerrors.Errorf("invalid license: %w", err)
		return
	}

	if cmd.Hardware, err = spaceprimitive.NewHardware(req.Hardware, req.SDK); err != nil {
		err = xerrors.Errorf("invalid hardware: %w", err)
		return
	}

	if cmd.Visibility, err = primitive.NewVisibility(req.Visibility); err != nil {
		err = xerrors.Errorf("invalid visibility: %w", err)
		return
	}

	if cmd.Fullname, err = primitive.NewMSDFullname(req.Fullname); err != nil {
		err = xerrors.Errorf("invalid fullname: %w", err)
		return
	}

	if cmd.AvatarId, err = primitive.NewAvatar(req.AvatarId); err != nil {
		err = xerrors.Errorf("invalid avatar id: %w", err)
		return
	}

	if cmd.BaseImage, err = spaceprimitive.NewBaseImage(req.BaseImage, req.Hardware); err != nil {
		err = xerrors.Errorf("invalid base image: %w", err)
		return
	}

	// always init readme
	cmd.InitReadme = true

	return
}

// reqToUpdateSpace
type reqToUpdateSpace struct {
	SDK        *string `json:"sdk"`
	Desc       *string `json:"desc"`
	AvatarId   *string `json:"space_avatar_id"`
	Fullname   *string `json:"fullname"`
	Hardware   *string `json:"hardware"`
	Visibility *string `json:"visibility"`
}

func (p *reqToUpdateSpace) action() (str string) {
	if p.Visibility != nil {
		str += fmt.Sprintf("visibility = %s", *p.Visibility)
	}

	return
}

func (p *reqToUpdateSpace) toCmd() (cmd app.CmdToUpdateSpace, err error) {
	if p.Desc != nil {
		if cmd.Desc, err = primitive.NewMSDDesc(*p.Desc); err != nil {
			return
		}
	}

	if p.Fullname != nil {
		if cmd.Fullname, err = primitive.NewMSDFullname(*p.Fullname); err != nil {
			return
		}
	}

	if p.Visibility != nil {
		if cmd.Visibility, err = primitive.NewVisibility(*p.Visibility); err != nil {
			return
		}
	}

	if p.AvatarId != nil {
		if cmd.AvatarId, err = primitive.NewAvatar(*p.AvatarId); err != nil {
			return
		}
	}

	return
}

// reqToDisableSpace
type reqToDisableSpace struct {
	Reason string `json:"reason"`
}

func (p *reqToDisableSpace) action() (str string) {
	str += fmt.Sprintf("reason = %s", p.Reason)

	return
}

func (p *reqToDisableSpace) toCmd() (cmd app.CmdToDisableSpace, err error) {
	cmd.Disable = true

	if cmd.DisableReason, err = primitive.NewDisableReason(p.Reason); err != nil {
		return
	}

	return
}

// reqToListUserSpaces
type reqToListUserSpaces struct {
	Name string `form:"name"`
	controller.CommonListRequest
}

func (req *reqToListUserSpaces) toCmd() (cmd app.CmdToListSpaces, err error) {
	cmd.Name = req.Name
	cmd.Count = req.Count

	if req.SortBy == "" {
		req.SortBy = primitive.SortByGlobal
	}
	if cmd.SortType, err = primitive.NewSortType(req.SortBy); err != nil {
		return
	}

	if v := req.CountPerPage; v <= 0 || v > config.MaxCountPerPage {
		cmd.CountPerPage = config.MaxCountPerPage
	} else {
		cmd.CountPerPage = v
	}

	if v := req.PageNum; v <= 0 {
		cmd.PageNum = firstPage
	} else {
		if v > (math.MaxInt / cmd.CountPerPage) {
			err = errors.New("invalid page num")

			return
		}
		cmd.PageNum = v
	}

	return
}

// reqToListGlobalSpaces
type reqToListGlobalSpaces struct {
	Task         string `form:"task"`
	License      string `form:"license"`
	Hardware     string `form:"hardware"`
	Framework    string `form:"framework"`
	HasAppFile   bool   `form:"has_app_file"`
	HardwareType string `form:"hardware_type"`
	reqToListUserSpaces
}

func (req *reqToListGlobalSpaces) toCmd() (app.CmdToListSpaces, error) {
	cmd, err := req.reqToListUserSpaces.toCmd()
	if err != nil {
		return cmd, err
	}

	if req.License != "" {
		if cmd.License, err = primitive.NewLicense(req.License); err != nil {
			return cmd, xerrors.Errorf("invalid license: %w", err)
		}
	}

	if req.Task != "" {
		if cmd.Labels.Task, err = spaceprimitive.NewTask(req.Task); err != nil {
			return cmd, xerrors.Errorf("invalid task: %w", err)
		}
	}

	if req.Hardware != "" {
		if !spaceprimitive.IsValidHardware(req.Hardware) {
			return cmd, xerrors.Errorf("invalid hardware: %s", req.Hardware)
		}
		cmd.Hardware = spaceprimitive.CreateHardware(req.Hardware)
	}

	if req.Framework != "" {
		if !spaceprimitive.IsValidFramework(req.Framework) {
			return cmd, xerrors.Errorf("invalid framework: %s, shoulde be %s or %s",
				req.Framework, spaceprimitive.PyTorch, spaceprimitive.MindSpore)
		}
		cmd.Framework = req.Framework
	}

	cmd.HasAppFile = req.HasAppFile
	cmd.Labels.HardwareType = req.HardwareType
	return cmd, nil
}

// restfulReqToListSpaces
type restfulReqToListSpaces struct {
	Owner string `form:"owner"`

	reqToListGlobalSpaces
}

func (req *restfulReqToListSpaces) toCmd() (app.CmdToListSpaces, error) {
	cmd, err := req.reqToListGlobalSpaces.toCmd()
	if err != nil {
		return cmd, err
	}

	if req.Owner != "" {
		if cmd.Owner, err = primitive.NewAccount(req.Owner); err != nil {
			return cmd, err
		}
	}

	return cmd, nil
}

type reqReportSpaceEmail struct {
	SpaceName string `json:"name"`
	Msg       string `json:"msg"`
	Owner     string `json:"owner"`
}

func (p *reqReportSpaceEmail) toCmd() (cmd app.CmdToReportDatasetEmail, err error) {
	if p.SpaceName == "" || p.Msg == "" {
		e := fmt.Errorf("missing parameters")
		err = allerror.NewNoPermission(e.Error(), e)
		return app.CmdToReportDatasetEmail{}, err
	}
	if cmd.SpaceName, err = primitive.NewMSDName(p.SpaceName); err != nil {
		return cmd, xerrors.Errorf("%w", err)
	}
	if cmd.Owner, err = primitive.NewAccount(p.Owner); err != nil {
		return cmd, xerrors.Errorf("%w", err)
	}
	cmd.Msg = p.Msg
	return cmd, nil
}

// spaceDetail
type spaceDetail struct {
	Liked         bool   `json:"liked"`
	OwnerAvatarId string `json:"avatar_id"`
	OwnerType     int    `json:"owner_type"`

	*app.SpaceDTO
}

// userSpacesInfo
type userSpacesInfo struct {
	Owner     string `json:"owner"`
	AvatarId  string `json:"avatar_id"`
	OwnerType int    `json:"owner_type"`

	*app.SpacesDTO
}

// spaceInfo
type spaceInfo struct {
	AvatarId  string `json:"avatar_id"`
	OwnerType int    `json:"owner_type"`
	Owner     string `json:"owner"`

	*repository.SpaceSummary
}

// spacesInfo
type spacesInfo struct {
	Total  int         `json:"total"`
	Spaces []spaceInfo `json:"spaces"`
}

// ModeIds is []string{"owner/name"}
type ModeIds struct {
	Ids []string `json:"ids"`
}

func (req *ModeIds) toCmd() []*domain.ModelIndex {
	modelsIndex := make([]*domain.ModelIndex, 0, len(req.Ids))

	for _, id := range req.Ids {
		index := strings.Split(id, "/")
		if len(index) != repoNameSplitedLen {
			logrus.Debugf("invalid model_id: %s", id)
			continue
		}
		owner, err := primitive.NewAccount(index[0])
		if err != nil {
			logrus.Debugf("invalid owner: %s", owner)
			continue
		}
		name, err := primitive.NewMSDName(index[1])
		if err != nil {
			logrus.Debugf("invalid model name: %s", name)
			continue
		}
		modelIndex := domain.ModelIndex{Owner: owner, Name: name}
		modelsIndex = append(modelsIndex, &modelIndex)
	}

	return modelsIndex
}

type reqToCreateSpaceVariable struct {
	Name  *string `json:"name"       required:"true"`
	Desc  *string `json:"desc"`
	Value *string `json:"value"`
}

func (p *reqToCreateSpaceVariable) toCmd() (cmd app.CmdToCreateSpaceVariable, err error) {
	if p.Name != nil {
		if cmd.Name, err = spaceprimitive.NewENVName(*p.Name); err != nil {
			err = xerrors.Errorf("failed to create env name, err:%w", err)
			return
		}
	}

	if p.Desc != nil {
		if cmd.Desc, err = primitive.NewMSDDesc(*p.Desc); err != nil {
			err = xerrors.Errorf("failed to create env desc, err:%w", err)
			return
		}
	}

	if p.Value != nil {
		if cmd.Value, err = spaceprimitive.NewENVValue(*p.Value); err != nil {
			err = xerrors.Errorf("failed to create env value, err:%w", err)
			return
		}
	}

	return
}

// reqToUpdateSpaceVariable
type reqToUpdateSpaceVariable struct {
	Value *string `json:"value"`
	Desc  *string `json:"desc"`
}

func (p *reqToUpdateSpaceVariable) action() (str string) {
	if p.Value != nil {
		str += fmt.Sprintf("value = %s", *p.Value)
	}

	if p.Desc != nil {
		str += fmt.Sprintf("desc = %s", *p.Desc)
	}

	return
}

func (p *reqToUpdateSpaceVariable) toCmd() (cmd app.CmdToUpdateSpaceVariable, err error) {
	if p.Desc != nil {
		if cmd.Desc, err = primitive.NewMSDDesc(*p.Desc); err != nil {
			return
		}
	}

	if p.Value != nil {
		if cmd.Value, err = spaceprimitive.NewENVValue(*p.Value); err != nil {
			return
		}
	}

	return
}

type reqToCreateSpaceSecret struct {
	Name  *string `json:"name"       required:"true"`
	Desc  *string `json:"desc"`
	Value *string `json:"value"`
}

func (p *reqToCreateSpaceSecret) toCmd() (cmd app.CmdToCreateSpaceSecret, err error) {
	if p.Name != nil {
		if cmd.Name, err = spaceprimitive.NewENVName(*p.Name); err != nil {
			err = xerrors.Errorf("failed to create env name, err:%w", err)
			return
		}
	}

	if p.Desc != nil {
		if cmd.Desc, err = primitive.NewMSDDesc(*p.Desc); err != nil {
			err = xerrors.Errorf("failed to create env desc, err:%w", err)
			return
		}
	}

	if p.Value != nil {
		if cmd.Value, err = spaceprimitive.NewENVValue(*p.Value); err != nil {
			err = xerrors.Errorf("failed to create env value, err:%w", err)
			return
		}
	}

	return
}

// reqToUpdateSpaceSecret
type reqToUpdateSpaceSecret struct {
	Value *string `json:"value"`
	Desc  *string `json:"desc"`
}

func (p *reqToUpdateSpaceSecret) toCmd() (cmd app.CmdToUpdateSpaceSecret, err error) {
	if p.Desc != nil {
		if cmd.Desc, err = primitive.NewMSDDesc(*p.Desc); err != nil {
			return
		}
	}

	if p.Value != nil {
		if cmd.Value, err = spaceprimitive.NewENVValue(*p.Value); err != nil {
			return
		}
	}

	return
}

type localCMD space.LocalCMD

func (req *localCMD) toCmd() string {
	return req.Cmd
}

type localEnvInfo space.LocalEnvInfo

func (req *localEnvInfo) toCmd() string {
	return req.EnvInfo
}

type spaceRecommendInfo struct {
	OwnerType int `json:"owner_type"`
	*repository.SpaceSummary
}

type spacesRecommendInfo struct {
	Spaces []spaceRecommendInfo `json:"spaces"`
}

type reqToResetLabel struct {
	Licenses []string
	Task     string
}

func (req *reqToResetLabel) toCmd() (cmd app.CmdToResetLabels, err error) {
	if len(req.Licenses) > 0 {
		cmd.Licenses = sets.New[string](req.Licenses...)
	}

	if cmd.Task, err = spaceprimitive.NewTask(req.Task); err != nil {
		return cmd, xerrors.Errorf("invalid task: %w", err)
	}

	return
}

// reqToNotifyUpdateCode
type reqToNotifyUpdateCode struct {
	CommitId string `json:"commit_id"`
	SdkType  string `json:"sdk_type"`
}

func (req *reqToNotifyUpdateCode) toCmd() (cmd app.CmdToNotifyUpdateCode, err error) {
	cmd.CommitId = req.CommitId
	if req.SdkType == "" {
		cmd.HasHtml = false
		cmd.HasApp = false
		return
	}
	if req.SdkType == spaceprimitive.StaticSdk.SDK() {
		cmd.HasHtml = true
		return
	}
	if req.SdkType == spaceprimitive.GradioSdk.SDK() {
		cmd.HasApp = true
		return
	}
	cmd.HasHtml = true
	cmd.HasApp = true
	return
}

func toUploadCoverCmd(ctx *gin.Context, user primitive.Account) (app.CmdToUploadCover, error) {
	f, err := ctx.FormFile("file")
	if err != nil {
		return app.CmdToUploadCover{}, xerrors.Errorf("failed to parse request param: %w", err)
	}

	if f.Size > config.MaxCoverFileSize {
		err = errors.New("file too big")

		return app.CmdToUploadCover{}, err
	}

	name, err := primitive.NewFileName(f.Filename)
	if err != nil {
		return app.CmdToUploadCover{}, err
	}

	if !name.IsPictureName() {
		err = errors.New("file format error")

		return app.CmdToUploadCover{}, err
	}

	p, err := f.Open()
	if err != nil {
		return app.CmdToUploadCover{}, xerrors.Errorf("can not get file: %w", err)
	}

	defer p.Close()

	cmd := app.CmdToUploadCover{
		User:  user,
		Image: p,
		FileName: commonctl.GetSaltHash(fmt.Sprintf("%s%s%v",
			user.Account(), f.Filename, utils.Now())) + name.GetFormat(),
	}

	return cmd, nil
}
