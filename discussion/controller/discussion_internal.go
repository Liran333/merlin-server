package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/discussion/app"
)

func AddRouterForDiscussionInternalController(
	r *gin.RouterGroup,
	m middleware.UserMiddleWare,
	s app.IssueInternalService,
) {

	ctl := DiscussionInternalController{
		app: s,
	}

	r.PUT("/v1/discussion/issue/:id", m.Write, ctl.UpdateCommentCount)
}

type DiscussionInternalController struct {
	app app.IssueInternalService
}

// @Summary  Update comment
// @Description  update comment count
// @Tags     DiscussionInternal
// @Param    id    path    string                     true    "id of issue"
// @Param    body  body    reqToUpdateCommentCount    true    "body of update comment count"
// @Success  202    {object}    commonctl.ResponseData{data=string,msg=string,code=string}
// @Router   /v1/discussion/issue/{id} [put]
func (ctl *DiscussionInternalController) UpdateCommentCount(ctx *gin.Context) {
	issueId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	req := reqToUpdateCommentCount{}
	if err = ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	if err = ctl.app.UpdateCommentCount(ctx.Request.Context(), issueId, req.Count); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}
