//nolint:typecheck
package controller

import (
	"errors"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	sdk "github.com/openmerlin/merlin-sdk/activityapp"
	"github.com/openmerlin/merlin-server/activity/app"
	"github.com/openmerlin/merlin-server/common/controller"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	modelapp "github.com/openmerlin/merlin-server/models/app"
	orgapp "github.com/openmerlin/merlin-server/organization/app"
	spaceapp "github.com/openmerlin/merlin-server/space/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

// AddRouteForActivityWebController adds a router for the ActivityWebController with the given middleware.
func AddRouteForActivityWebController(
	r *gin.RouterGroup,
	s app.ActivityAppService,
	m middleware.UserMiddleWare,
	o orgapp.OrgService,
	u userapp.UserService,
	d modelapp.ModelAppService,
	p spaceapp.SpaceAppService,
	rl middleware.RateLimiter,
) {
	ctl := ActivityWebController{
		ActivityController: ActivityController{
			appService:     s,
			userMiddleWare: m,
			user:           u,
			org:            o,
			model:          d,
			space:          p,
		},
	}

	r.GET("/v1/user/activity", m.Read, rl.CheckLimit, ctl.List)
	r.POST("/v1/like", m.Write, ctl.Add)
	r.DELETE("/v1/like", m.Write, ctl.Delete)
}

func (req *reqToListUserActivities) toCmd() (cmd app.CmdToListActivities, err error) {
	cmd.Count = req.Count
	cmd.Model = req.Model
	cmd.Space = req.Space
	cmd.Like = req.Like
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

// ActivityWebController is a struct that holds the app service for model web operations.
type ActivityWebController struct {
	ActivityController
}

// reqToListUserModels
type reqToListUserActivities struct {
	Model string `form:"model"`
	Space string `form:"space"`
	Like  string `form:"like"`
	controller.CommonListRequest
}

// @Summary  List
// @Description  get activities
// @Tags     Activity
// @Param    space     query  string  false "filter by space"
// @Param    model  query  string  false "filter by model"
// @Param    like  query  string  false "filter by like"
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData
// @Failure  400  {object}  commonctl.ResponseData
// @Router /web/v1/user/activity [get]
func (ctl *ActivityWebController) List(ctx *gin.Context) {
	// Bind query parameters to request struct
	var req reqToListUserActivities
	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	// Convert request to command
	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	// Get user from middleware
	user := ctl.userMiddleWare.GetUser(ctx)

	// Find organizations visible to the user
	orgNames, _ := ctl.org.GetByUser(user, user)

	// Prepare list of names including the user's account name
	var list []primitive.Account

	list = append(list, user)
	for _, org := range orgNames {
		acc := primitive.CreateAccount(org.Name)
		if err != nil {
			logrus.Errorf("Error creating account for org %s: %v", org.Name, err)
			continue
		}
		list = append(list, acc)
	}

	// List activities based on the prepared list and command
	dto, err := ctl.appService.List(list, &cmd)
	if err != nil {
		commonctl.SendError(ctx, err)
		return
	}

	commonctl.SendRespOfGet(ctx, &dto)
}

// @Summary  Add
// @Description  add a like record in the activity table
// @Tags     Activity
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData
// @Failure  400  {object}  commonctl.ResponseData
// @Router /web/v1/like [post]
func (ctl *ActivityWebController) Add(ctx *gin.Context) {
	var req sdk.ReqToCreateActivity

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	req.Owner = user.Account()

	cmd, err := ConvertReqToCreateActivityToCmd(&req)

	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	//_, err = ctl.model.AddLike(cmd.Resource.Index)
	//
	//if err != nil {
	//	return
	//}
	//
	//_, err = ctl.space.AddLike(primitive.CreateIdentity(cmd.ResourceId))
	//
	//if err != nil {
	//	return
	//}

	if err := ctl.appService.Create(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  Delete
// @Description  Delete a like record in the activity table
// @Tags     Activity
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData
// @Failure  400  {object}  commonctl.ResponseData
// @Router /web/v1/like [delete]
func (ctl *ActivityWebController) Delete(ctx *gin.Context) {
	var req sdk.ReqToDeleteActivity

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)
		return
	}

	cmd, err := ConvertReqToDeleteActivityToCmd(&req)

	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	if err := ctl.appService.Delete(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}
