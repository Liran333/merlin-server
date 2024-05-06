/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/openmerlin/merlin-sdk/models"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/app"
	"github.com/openmerlin/merlin-server/models/domain/repository"
)

const (
	firstPage    = 1
	labelSpliter = ","
)

type reqToCreateModel struct {
	Name       string `json:"name"       required:"true"`
	Desc       string `json:"desc"`
	Owner      string `json:"owner"      required:"true"`
	License    string `json:"license"    required:"true"`
	Fullname   string `json:"fullname"`
	Visibility string `json:"visibility" required:"true"`
}

func (req *reqToCreateModel) action() string {
	return fmt.Sprintf("create model of %s/%s", req.Owner, req.Name)
}

func (req *reqToCreateModel) toCmd() (cmd app.CmdToCreateModel, err error) {
	if cmd.Name, err = primitive.NewMSDName(req.Name); err != nil {
		return
	}

	if cmd.Desc, err = primitive.NewMSDDesc(req.Desc); err != nil {
		return
	}

	if cmd.Owner, err = primitive.NewAccount(req.Owner); err != nil {
		return
	}

	if cmd.License, err = primitive.NewLicense(req.License); err != nil {
		return
	}

	if cmd.Visibility, err = primitive.NewVisibility(req.Visibility); err != nil {
		return
	}

	if cmd.Fullname, err = primitive.NewMSDFullname(req.Fullname); err != nil {
		return
	}

	// always init readme
	cmd.InitReadme = true

	return
}

// reqToUpdateModel
type reqToUpdateModel struct {
	Desc       *string `json:"desc"`
	Fullname   *string `json:"fullname"`
	Visibility *string `json:"visibility"`
}

func (p *reqToUpdateModel) action() (str string) {
	if p.Visibility != nil {
		str += fmt.Sprintf("visibility = %s", *p.Visibility)
	}

	return
}

func (p *reqToUpdateModel) toCmd() (cmd app.CmdToUpdateModel, err error) {
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

	return
}

// reqToDisableModel
type reqToDisableModel struct {
	Reason string `json:"reason"`
}

func (p *reqToDisableModel) action() (str string) {
	str += fmt.Sprintf("reason = %s", p.Reason)

	return
}

func (p *reqToDisableModel) toCmd() (cmd app.CmdToDisableModel, err error) {
	cmd.Disable = true

	if cmd.DisableReason, err = primitive.NewDisableReason(p.Reason); err != nil {
		return
	}

	return
}

// reqToListUserModels
type reqToListUserModels struct {
	Name string `form:"name"`
	controller.CommonListRequest
}

func (req *reqToListUserModels) toCmd() (cmd app.CmdToListModels, err error) {
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

// reqToListGlobalModels
type reqToListGlobalModels struct {
	Task       string `form:"task"`
	License    string `form:"license"`
	Frameworks string `form:"frameworks"`

	reqToListUserModels
}

func (req *reqToListGlobalModels) toCmd() (app.CmdToListModels, error) {
	cmd, err := req.reqToListUserModels.toCmd()
	if err != nil {
		return cmd, err
	}

	// TODO check each label if it is valid

	cmd.Labels.Task = req.Task

	if req.License != "" {
		if cmd.License, err = primitive.NewLicense(req.License); err != nil {
			return cmd, err
		}
	}

	cmd.Labels.Frameworks = toStringsSets(req.Frameworks)

	return cmd, nil
}

func toStringsSets(v string) sets.Set[string] {
	if v == "" {
		return nil
	}

	items := strings.Split(v, labelSpliter)

	return sets.New[string](items...)
}

// restfulReqToListModels
type restfulReqToListModels struct {
	Owner string `form:"owner"`

	reqToListGlobalModels
}

func (req *restfulReqToListModels) toCmd() (app.CmdToListModels, error) {
	cmd, err := req.reqToListGlobalModels.toCmd()
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

// modelDetail
type modelDetail struct {
	Liked    bool   `json:"liked"`
	AvatarId string `json:"avatar_id"`

	*app.ModelDTO
}

// modelsInfo
type userModelsInfo struct {
	Owner     string `json:"owner"`
	AvatarId  string `json:"avatar_id"`
	OwnerType int    `json:"owner_type"`

	*app.ModelsDTO
}

// modelsInfo
type modelInfo struct {
	AvatarId  string `json:"avatar_id"`
	OwnerType int    `json:"owner_type"`
	Owner     string `json:"owner"`

	*repository.ModelSummary
}

// modelsInfo
type modelsInfo struct {
	Total  int         `json:"total"`
	Models []modelInfo `json:"models"`
}

// models statistics
type modelStatistics struct {
	DownloadCount int `json:"download_count"`
}

func (s *modelStatistics) toCmd() app.CmdToUpdateStatistics {
	return app.CmdToUpdateStatistics{
		DownloadCount: s.DownloadCount,
	}
}

type useInOpenmind models.UseInOpenmind

func (req *useInOpenmind) toCmd() string {
	return req.UseInOpenmind
}
