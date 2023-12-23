package controller

import (
	"errors"
	"math"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"

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
	InitReadme bool   `json:"init_readme"`
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

	cmd.InitReadme = req.InitReadme

	return
}

// reqToUpdateModel
type reqToUpdateModel struct {
	Name       *string `json:"name"`
	Desc       *string `json:"desc"`
	Fullname   *string `json:"fullname"`
	Visibility *string `json:"visibility"`
}

func (p *reqToUpdateModel) toCmd() (cmd app.CmdToUpdateModel, err error) {
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

// reqToListUserModels
type reqToListUserModels struct {
	Name         string `form:"name"`
	SortBy       string `form:"sort_by"`
	Count        bool   `form:"count"`
	PageNum      int    `form:"page_num"`
	CountPerPage int    `form:"count_per_page"`
}

func (req *reqToListUserModels) toCmd() (cmd app.CmdToListModels, err error) {
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

// reqToListGlobalModels
type reqToListGlobalModels struct {
	Task       string `form:"task"`
	Others     string `form:"others"`
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

	if v := strings.Split(req.Others, labelSpliter); len(v) > 0 {
		cmd.Labels.Others = sets.Set[string](sets.NewString(v...))
	}

	if v := strings.Split(req.Frameworks, labelSpliter); len(v) > 0 {
		cmd.Labels.Frameworks = sets.Set[string](sets.NewString(v...))
	}

	return cmd, nil
}

// restfulReqToListModels
type restfulReqToListModels struct {
	owner string `form:"owner"`

	reqToListGlobalModels
}

func (req *restfulReqToListModels) toCmd() (app.CmdToListModels, error) {
	cmd, err := req.reqToListGlobalModels.toCmd()
	if err != nil {
		return cmd, err
	}

	cmd.Owner, err = primitive.NewAccount(req.owner)

	return cmd, err
}

// modelDetail
type modelDetail struct {
	Liked    bool   `json:"liked"`
	AvatarId string `json:"avatar_id"`

	*app.ModelDTO
}

// modelsInfo
type userModelsInfo struct {
	Owner    string `json:"owner"`
	AvatarId string `json:"avatar_id"`

	*app.ModelsDTO
}

// modelsInfo
type modelInfo struct {
	AvatarId string `json:"avatar_id"`

	*repository.ModelSummary
}

// modelsInfo
type modelsInfo struct {
	Total  int         `json:"total"`
	Models []modelInfo `json:"models"`
}
