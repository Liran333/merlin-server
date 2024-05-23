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

	// "github.com/openmerlin/merlin-sdk/datasets"

	"github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/datasets/app"
	"github.com/openmerlin/merlin-server/datasets/domain/repository"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	firstPage    = 1
	labelSpliter = ","
)

type reqToCreateDataset struct {
	Name       string `json:"name"       required:"true"`
	Desc       string `json:"desc"`
	Owner      string `json:"owner"      required:"true"`
	License    string `json:"license"    required:"true"`
	Fullname   string `json:"fullname"`
	Visibility string `json:"visibility" required:"true"`
}

func (req *reqToCreateDataset) action() string {
	return fmt.Sprintf("create dataset of %s/%s", req.Owner, req.Name)
}

func (req *reqToCreateDataset) toCmd() (cmd app.CmdToCreateDataset, err error) {
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

// reqToUpdateDataset
type reqToUpdateDataset struct {
	Desc       *string `json:"desc"`
	Fullname   *string `json:"fullname"`
	Visibility *string `json:"visibility"`
}

func (p *reqToUpdateDataset) action() (str string) {
	if p.Visibility != nil {
		str += fmt.Sprintf("visibility = %s", *p.Visibility)
	}

	return
}

func (p *reqToUpdateDataset) toCmd() (cmd app.CmdToUpdateDataset, err error) {
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

// reqToDisableDataset
type reqToDisableDataset struct {
	Reason string `json:"reason"`
}

func (p *reqToDisableDataset) action() (str string) {
	str += fmt.Sprintf("reason = %s", p.Reason)

	return
}

func (p *reqToDisableDataset) toCmd() (cmd app.CmdToDisableDataset, err error) {
	cmd.Disable = true

	if cmd.DisableReason, err = primitive.NewDisableReason(p.Reason); err != nil {
		return
	}

	return
}

// reqToListUserDatasets
type reqToListUserDatasets struct {
	Name string `form:"name"`
	controller.CommonListRequest
}

func (req *reqToListUserDatasets) toCmd() (cmd app.CmdToListDatasets, err error) {
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

// reqToListGlobalDatasets
type reqToListGlobalDatasets struct {
	Task     string `form:"task"`
	License  string `form:"license"`
	Size     string `form:"size"`
	Language string `form:"language"`
	Domain   string `form:"domain"`

	reqToListUserDatasets
}

func (req *reqToListGlobalDatasets) toCmd() (app.CmdToListDatasets, error) {
	cmd, err := req.reqToListUserDatasets.toCmd()
	if err != nil {
		return cmd, err
	}

	// TODO check each label if it is valid
	cmd.Labels.Task = toStringsSets(req.Task)

	if req.License != "" {
		if cmd.License, err = primitive.NewLicense(req.License); err != nil {
			return cmd, err
		}
	}

	cmd.Labels.Size = req.Size

	cmd.Labels.Language = toStringsSets(req.Language)
	cmd.Labels.Domain = toStringsSets(req.Domain)

	return cmd, nil
}

func toStringsSets(v string) sets.Set[string] {
	if v == "" {
		return nil
	}

	items := strings.Split(v, labelSpliter)

	return sets.New[string](items...)
}

// restfulReqToListDatasets
type restfulReqToListDatasets struct {
	Owner string `form:"owner"`

	reqToListGlobalDatasets
}

func (req *restfulReqToListDatasets) toCmd() (app.CmdToListDatasets, error) {
	cmd, err := req.reqToListGlobalDatasets.toCmd()
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

type datasetDetail struct {
	Liked    bool   `json:"liked"`
	AvatarId string `json:"avatar_id"`

	*app.DatasetDTO
}

type userDatasetsInfo struct {
	Owner     string `json:"owner"`
	AvatarId  string `json:"avatar_id"`
	OwnerType int    `json:"owner_type"`

	*app.DatasetsDTO
}

type datasetInfo struct {
	AvatarId  string `json:"avatar_id"`
	OwnerType int    `json:"owner_type"`
	Owner     string `json:"owner"`

	*repository.DatasetSummary
}

type datasetsInfo struct {
	Total    int           `json:"total"`
	Datasets []datasetInfo `json:"datasets"`
}

type datasetStatistics struct {
	DownloadCount int `json:"download_count"`
}

func (s *datasetStatistics) toCmd() app.CmdToUpdateStatistics {
	return app.CmdToUpdateStatistics{
		DownloadCount: s.DownloadCount,
	}
}
