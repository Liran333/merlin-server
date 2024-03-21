/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/domain"
	"github.com/openmerlin/merlin-server/space/app"
	spaceprimitive "github.com/openmerlin/merlin-server/space/domain/primitive"
	"github.com/openmerlin/merlin-server/space/domain/repository"
	"github.com/sirupsen/logrus"
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
	Fullname   string `json:"fullname"`
	Visibility string `json:"visibility" required:"true"`
	InitReadme bool   `json:"init_readme"`
}

func (req *reqToCreateSpace) action() string {
	return fmt.Sprintf("create space of %s/%s", req.Owner, req.Name)
}

func (req *reqToCreateSpace) toCmd() (cmd app.CmdToCreateSpace, err error) {
	if cmd.SDK, err = spaceprimitive.NewSDK(req.SDK); err != nil {
		return
	}

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

	if cmd.Hardware, err = spaceprimitive.NewHardware(req.Hardware, req.SDK); err != nil {
		return
	}

	if cmd.Visibility, err = primitive.NewVisibility(req.Visibility); err != nil {
		return
	}

	if cmd.Fullname, err = primitive.NewMSDFullname(req.Fullname); err != nil {
		return
	}

	cmd.InitReadme = req.InitReadme

	return
}

// reqToUpdateSpace
type reqToUpdateSpace struct {
	SDK        *string `json:"sdk"`
	Name       *string `json:"name"`
	Desc       *string `json:"desc"`
	Fullname   *string `json:"fullname"`
	Hardware   *string `json:"hardware"`
	Visibility *string `json:"visibility"`
}

func (p *reqToUpdateSpace) action() (str string) {
	if p.Name != nil {
		str += fmt.Sprintf("name = %s", *p.Name)
	}

	if p.Visibility != nil {
		str += fmt.Sprintf("visibility = %s", *p.Visibility)
	}

	return
}

func (p *reqToUpdateSpace) toCmd() (cmd app.CmdToUpdateSpace, err error) {
	if p.Name != nil {
		if cmd.Name, err = primitive.NewMSDName(*p.Name); err != nil {
			return
		}
	}

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

// reqToListUserSpaces
type reqToListUserSpaces struct {
	Name string `form:"name"`
	controller.CommonListRequest
}

func (req *reqToListUserSpaces) toCmd() (cmd app.CmdToListSpaces, err error) {
	cmd.Name = req.Name
	cmd.Count = req.Count

	if req.SortBy == "" {
		cmd.SortType = primitive.SortTypeRecentlyUpdated
	} else {
		if cmd.SortType, err = primitive.NewSortType(req.SortBy); err != nil {
			return
		}
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
	Task       string `form:"task"`
	Others     string `form:"others"`
	License    string `form:"license"`
	Frameworks string `form:"frameworks"`

	reqToListUserSpaces
}

func (req *reqToListGlobalSpaces) toCmd() (app.CmdToListSpaces, error) {
	cmd, err := req.reqToListUserSpaces.toCmd()
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

	cmd.Labels.Others = toStringsSets(req.Others)
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

// spaceDetail
type spaceDetail struct {
	Liked    bool   `json:"liked"`
	AvatarId string `json:"avatar_id"`

	*app.SpaceDTO
}

// userSpacesInfo
type userSpacesInfo struct {
	Owner    string `json:"owner"`
	AvatarId string `json:"avatar_id"`

	*app.SpacesDTO
}

// spaceInfo
type spaceInfo struct {
	AvatarId string `json:"avatar_id"`

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
		if cmd.Name, err = primitive.NewMSDName(*p.Name); err != nil {
			return
		}
	}

	if p.Desc != nil {
		if cmd.Desc, err = primitive.NewMSDDesc(*p.Desc); err != nil {
			return
		}
	}

	if p.Value != nil {
		if cmd.Value, err = primitive.NewMSDName(*p.Value); err != nil {
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
		if cmd.Value, err = primitive.NewMSDName(*p.Value); err != nil {
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
		if cmd.Name, err = primitive.NewMSDName(*p.Name); err != nil {
			return
		}
	}

	if p.Desc != nil {
		if cmd.Desc, err = primitive.NewMSDDesc(*p.Desc); err != nil {
			return
		}
	}

	if p.Value != nil {
		if cmd.Value, err = primitive.NewMSDName(*p.Value); err != nil {
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
		if cmd.Value, err = primitive.NewMSDName(*p.Value); err != nil {
			return
		}
	}

	return
}