package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/app"
)

func AddRouteForSoftwarePkgController(
	r *gin.RouterGroup,
	s app.ModelAppService,
) {
	ctl := ModelController{
		appService: s,
	}

	must := middleware.UserMiddleware().Must
	optional := middleware.UserMiddleware().Optional

	r.POST(`/v1/model`, must, ctl.Create)
	r.DELETE("/v1/model/:id", must, ctl.Delete)
	r.PUT("/v1/model/:id", must, ctl.Update)
	r.GET("/v1/model/:owner/:name", optional, ctl.Get)
	r.GET("/v1/model/:owner", optional, ctl.List)
	r.GET("/v1/model/:owner", optional, ctl.ListGlobal)
}

type ModelController struct {
	appService app.ModelAppService
}

// @Summary  Create
// @Description  create model
// @Tags     Model
// @Param    body  body      reqToCreateModel  true  "body of creating model"
// @Accept   json
// @Success  201   {object}  commonctl.ResponseData
// @Router   /v1/model [post]
func (ctl *ModelController) Create(ctx *gin.Context) {
	req := reqToCreateModel{}
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := middleware.UserMiddleware().GetUser(ctx)

	if v, err := ctl.appService.Create(user, &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, v)
	}
}

// @Summary  Delete
// @Description  delete model
// @Tags     Model
// @Param    id    path  string        true  "id of model"
// @Accept   json
// @Success  204
// @Router   /v1/model/{id} [delete]
func (ctl *ModelController) Delete(ctx *gin.Context) {
	user := middleware.UserMiddleware().GetUser(ctx)

	modelId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.appService.Delete(user, modelId); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  Update
// @Description  update model
// @Tags     Model
// @Param    id    path  string            true  "id of model"
// @Param    body  body  reqToUpdateModel  true  "body of updating model"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData
// @Router   /v1/model/{id} [put]
func (ctl *ModelController) Update(ctx *gin.Context) {
	req := reqToUpdateModel{}
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	modelId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	err = ctl.appService.Update(
		middleware.UserMiddleware().GetUser(ctx),
		modelId, &cmd,
	)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx)
	}
}

// @Summary  Get
// @Description  get model
// @Tags     Model
// @Param    owner  path  string  true  "owner of model"
// @Param    name   path  string  true  "name of model"
// @Accept   json
// @Success  200  {object}  modelDetail
// @Router   /v1/model/{owner}/{name} [get]
func (ctl *ModelController) Get(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	user := middleware.UserMiddleware().GetUser(ctx)

	dto, err := ctl.appService.GetByName(user, &index)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	detail := modelDetail{
		Liked:    true,
		ModelDTO: &dto,
	}

	if user != nil {
		//TODO check user like the model
	}

	/*
		TODO get avatar of owner

		avatar, err := ctl.user.GetUserAvatarId(owner)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}

		if avatar != nil {
			detail.AvatarId = avatar.AvatarId()
		}
	*/

	commonctl.SendRespOfGet(ctx, &detail)
}

// @Summary  List
// @Description  list model
// @Tags     Model
// @Param    owner           path   string  true   "owner of model"
// @Param    name            query  string  false  "name of model"
// @Param    count           query  bool    false  "whether to calculate the total"
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created"
// @Param    page_num        query  int     false  "page num which starts from 1"
// @Param    count_per_page  query  int     false  "count per page"
// @Accept   json
// @Success  200  {object}  modelsInfo
// @Router   /v1/model/{owner} [get]
func (ctl *ModelController) List(ctx *gin.Context) {
	var req reqToListUserModels
	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd.Owner, err = primitive.NewAccount(ctx.Param("owner"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := middleware.UserMiddleware().GetUser(ctx)

	dto, err := ctl.appService.List(user, &cmd)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	result := modelsInfo{
		Owner:     cmd.Owner.Account(),
		ModelsDTO: &dto,
	}

	/*
		avatar, err := ctl.user.GetUserAvatarId(owner)
		if err != nil {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))

			return
		}

		if avatar != nil {
			result.AvatarId = avatar.AvatarId()
		}
	*/

	commonctl.SendRespOfGet(ctx, &result)
}

// @Summary  ListGlobal
// @Description  list global public model
// @Tags     Model
// @Param    name            query  string  false  "name of model"
// @Param    count           query  bool    false  "whether to calculate the total"
// @Param    labels          query  string  false  "labels, separate multiple each ones with commas"
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created"
// @Param    page_num        query  int     false  "page num which starts from 1"
// @Param    count_per_page  query  int     false  "count per page"
// @Accept   json
// @Success  200  {object}  app.ModelsDTO
// @Router   /v1/model [get]
func (ctl *ModelController) ListGlobal(ctx *gin.Context) {
	var req reqToListGlobalModels
	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := middleware.UserMiddleware().GetUser(ctx)

	if result, err := ctl.appService.List(user, &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, result)
	}
	// TODO: each model should include owner's avatar
}

func (ctl *ModelController) parseIndex(ctx *gin.Context) (index app.ModelIndex, err error) {
	index.Owner, err = primitive.NewAccount(ctx.Param("owner"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	index.Name, err = primitive.NewMSDName(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
	}

	return
}
