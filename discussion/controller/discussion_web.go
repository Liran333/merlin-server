package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/discussion/app"
)

func AddRouterForDiscussionWebController(
	r *gin.RouterGroup,
	m middleware.UserMiddleWare,
	l middleware.OperationLog,
	i app.IssueService,
	c app.CommentService,
	d app.DiscussionService,
) {
	ctl := DiscussionWebController{
		userMiddleWare:    m,
		issueService:      i,
		commentService:    c,
		discussionService: d,
	}

	r.POST("/v1/discussion/:resource_id/issue", m.Write, l.Write, ctl.CreateIssue)
	r.GET("/v1/discussion/:resource_id/issue/:id", m.Optional, ctl.GetIssue)
	r.GET("/v1/discussion/:resource_id/issue", m.Optional, ctl.ListIssue)
	r.GET("/v1/discussion/:resource_id/issue/count", m.Optional, ctl.ListIssuesCount)
	r.PUT("/v1/discussion/:resource_id/issue/:id/close", m.Write, l.Write, ctl.CloseIssue)
	r.PUT("/v1/discussion/:resource_id/issue/:id/reopen", m.Write, l.Write, ctl.ReopenIssue)
	r.PUT("/v1/discussion/:resource_id/close", m.Write, l.Write, ctl.CloseDiscussion)
	r.PUT("/v1/discussion/:resource_id/open", m.Write, l.Write, ctl.OpenDiscussion)

	r.POST("/v1/discussion/:resource_id/comment", m.Write, l.Write, ctl.CreateComment)
	r.PUT("/v1/discussion/:resource_id/comment/:id", m.Write, l.Write, ctl.UpdateComment)
	r.DELETE("/v1/discussion/:resource_id/comment/:id", m.Write, l.Write, ctl.DeleteComment)
	r.POST("/v1/discussion/:resource_id/comment/report/:id", m.Write, l.Write, ctl.ReportComment)
}

type DiscussionWebController struct {
	userMiddleWare    middleware.UserMiddleWare
	issueService      app.IssueService
	commentService    app.CommentService
	discussionService app.DiscussionService
}

// @Summary  Create issue
// @Description  create issue
// @Tags     DiscussionWeb
// @Param    resource_id    path    string true "id of model/space/datasets"
// @Param    body           body    reqToCreateIssue  true  "body of creating issue"
// @Accept   json
// @Security Bearer
// @Success  201    {object}    commonctl.ResponseData{data=string,msg=string,code=string}
// @Router   /v1/discussion/{resource_id}/issue [post]
func (ctl *DiscussionWebController) CreateIssue(ctx *gin.Context) {
	middleware.SetAction(ctx, "create issue")

	req := reqToCreateIssue{}
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())

	user := ctl.userMiddleWare.GetUser(ctx)
	cmd, err := req.toCreateIssueCmd(ctx.Param("resource_id"), user)
	if err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	if err = ctl.issueService.CreateIssue(ctx.Request.Context(), cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, nil)
	}
}

// @Summary  Get
// @Description  get issue
// @Tags     DiscussionWeb
// @Param    resource_id       path     string    true     "id of model/space/datasets"
// @Param    id                path     string    true     "id of issue"
// @Param    page_num          query    int       false    "page num which starts from 1" Mininum(1)
// @Param    count_per_page    query    int       false    "count per page" MaxCountPerPage(100)
// @Accept   json
// @Success  200    {object}    commonctl.ResponseData{data=IssueDetailDTO,msg=string,code=string}
// @Router   /v1/discussion/{resource_id}/issue/{id} [get]
func (ctl *DiscussionWebController) GetIssue(ctx *gin.Context) {
	issueId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	var req reqToGetIssue
	if err = ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)
	cmd, err := req.toGetIssueCmd(user, ctx.Param("resource_id"), issueId)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	issue, err := ctl.issueService.GetIssue(ctx.Request.Context(), cmd)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, &issue)
	}
}

// @Summary  Close issue
// @Description  close issue
// @Tags     DiscussionWeb
// @Param    resource_id     path     string     true     "id of model/space/datasets"
// @Param    id              path     string     true     "id of issue"
// @Success  202    {object}    commonctl.ResponseData{data=string,msg=string,code=string}
// @Router   /v1/discussion/{resource_id}/issue/{id}/close [put]
func (ctl *DiscussionWebController) CloseIssue(ctx *gin.Context) {
	issueId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx, fmt.Sprintf("close issue %d", issueId))

	user := ctl.userMiddleWare.GetUser(ctx)

	cmd, err := toCloseIssueCmd(user, ctx.Param("resource_id"), issueId)
	if err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	if err = ctl.issueService.CloseIssue(ctx.Request.Context(), cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  Reopen issue
// @Description  reopen issue
// @Tags     DiscussionWeb
// @Param    resource_id    path    string    true    "id of model/space/datasets"
// @Param    id             path    string    true    "id of issue"
// @Success  202    {object}    commonctl.ResponseData{data=string,msg=string,code=string}
// @Router   /v1/discussion/{resource_id}/issue/{id}/reopen [put]
func (ctl *DiscussionWebController) ReopenIssue(ctx *gin.Context) {
	issueId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx, fmt.Sprintf("reopen issue %d", issueId))

	user := ctl.userMiddleWare.GetUser(ctx)

	cmd, err := toCloseIssueCmd(user, ctx.Param("resource_id"), issueId)
	if err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	if err = ctl.issueService.ReopenIssue(ctx.Request.Context(), cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  Close discussion
// @Description  close discussion
// @Tags     DiscussionWeb
// @Param    resource_id    path    string    true    "id of model/space/datasets"
// @Success  202    {object}    commonctl.ResponseData{data=string,msg=string,code=string}
// @Router   /v1/discussion/{resource_id}/close [put]
func (ctl *DiscussionWebController) CloseDiscussion(ctx *gin.Context) {
	id, err := primitive.NewIdentity(ctx.Param("resource_id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx, fmt.Sprintf("close discussion %s", id.Identity()))

	user := ctl.userMiddleWare.GetUser(ctx)

	if err = ctl.discussionService.CloseDiscussion(ctx.Request.Context(), id, user); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  Open discussion
// @Description  open discussion
// @Tags     DiscussionWeb
// @Param    resource_id    path    string    true    "id of model/space/datasets"
// @Success  202    {object}    commonctl.ResponseData{data=string,msg=string,code=string}
// @Router   /v1/discussion/{resource_id}/open [put]
func (ctl *DiscussionWebController) OpenDiscussion(ctx *gin.Context) {
	id, err := primitive.NewIdentity(ctx.Param("resource_id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx, fmt.Sprintf("opoen discussion %s", id.Identity()))

	user := ctl.userMiddleWare.GetUser(ctx)

	if err = ctl.discussionService.ReopenDiscussion(ctx.Request.Context(), id, user); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  Issue count
// @Description  get issue count
// @Tags     DiscussionWeb
// @Param    resource_id    path    string    true    "id of model/space/datasets"
// @Success  200    {object}    commonctl.ResponseData{data=ListIssuesCountDTO,msg=string,code=string}
// @Router   /v1/discussion/{resource_id}/issue/count [get]
func (ctl *DiscussionWebController) ListIssuesCount(ctx *gin.Context) {
	id, err := primitive.NewIdentity(ctx.Param("resource_id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	data, err := ctl.issueService.ListIssuesCount(ctx.Request.Context(), user, id)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, &data)
	}
}

// @Summary  List
// @Description  list issues
// @Tags     DiscussionWeb
// @Param    resource_id       path     string    true     "id of model/space/datasets"
// @Param    page_num          query    int       false    "page num which starts from 1" Mininum(1)
// @Param    count_per_page    query    int       false    "count per page" MaxCountPerPage(100)
// @Param    status            query    string    false    "status of issue"
// @Accept   json
// @Success  200    {object}    commonctl.ResponseData{data=ListIssuesDTO,msg=string,code=string}
// @Router   /v1/discussion/{resource_id}/issue [get]
func (ctl *DiscussionWebController) ListIssue(ctx *gin.Context) {
	var req reqToListIssue
	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd, err := req.toListIssuesCmd(ctx.Param("resource_id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	if data, err := ctl.issueService.ListIssues(ctx.Request.Context(), user, cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, &data)
	}
}

// @Summary  Create comment
// @Description  create comment
// @Tags     DiscussionWeb
// @Param    resource_id    path    string                true    "id of model/space/datasets"
// @Param    body           body    reqToCreateComment    true    "body of creating comment"
// @Accept   json
// @Security Bearer
// @Success  201    {object}    commonctl.ResponseData{data=ItemDTO,msg=string,code=string}
// @Router   /v1/discussion/{resource_id}/comment [post]
func (ctl *DiscussionWebController) CreateComment(ctx *gin.Context) {
	middleware.SetAction(ctx, "create comment")

	var req reqToCreateComment
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)
	cmd, err := req.toCreateCommentCmd(user, ctx.Param("resource_id"))
	if err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	if dto, err := ctl.commentService.CreateIssueComment(ctx.Request.Context(), cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, &dto)
	}
}

// @Summary  Update comment
// @Description  update comment
// @Tags     DiscussionWeb
// @Param    resource_id     path     string              true    "id of model/space/datasets"
// @Param    id              path     string              true    "id of comment"
// @Param    body            body     reqToUpdateComment  true    "body of update comment"
// @Success  202    {object}    commonctl.ResponseData{data=string,msg=string,code=string}
// @Router   /v1/discussion/{resource_id}/comment/{id} [put]
func (ctl *DiscussionWebController) UpdateComment(ctx *gin.Context) {
	middleware.SetAction(ctx, "update comment")

	var req reqToUpdateComment
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	commentId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx, fmt.Sprintf("update comment %d", commentId))

	user := ctl.userMiddleWare.GetUser(ctx)
	cmd, err := req.toUpdateCommentCmd(commentId, user, ctx.Param("resource_id"))
	if err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	if err = ctl.commentService.UpdateIssueComment(ctx.Request.Context(), cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  Delete comment
// @Description  delete comment
// @Tags     DiscussionWeb
// @Param    resource_id    path    string    true    "id of model/space/datasets"
// @Param    id             path    string    true    "id of comment"
// @Success  204    {object}    commonctl.ResponseData{data=string,msg=string,code=string}
// @Router   /v1/discussion/{resource_id}/comment/{id} [delete]
func (ctl *DiscussionWebController) DeleteComment(ctx *gin.Context) {
	middleware.SetAction(ctx, "delete comment")

	commentId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx, fmt.Sprintf("delelte comment %d", commentId))

	user := ctl.userMiddleWare.GetUser(ctx)

	cmd, err := toDeleteCommentCmd(commentId, user, ctx.Param("resource_id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err = ctl.commentService.DeleteIssueComment(ctx.Request.Context(), cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  Report
// @Description  report comment
// @Tags     DiscussionWeb
// @Param    resource_id    path    string                true    "id of model/space/datasets"
// @Param    id             path    string                true    "id of comment"
// @Param    body           body    reqToReportComment    true    "body of report comment"
// @Accept   json
// @Security Bearer
// @Success  201    {object}    commonctl.ResponseData{data=string,msg=string,code=string}
// @Router   /v1/discussion/{resource_id}/comment/report/{id} [post]
func (ctl *DiscussionWebController) ReportComment(ctx *gin.Context) {
	commentId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx, fmt.Sprintf("report comment %d", commentId))

	var req reqToReportComment
	if err = ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)
	cmd, err1 := req.toReportCommentCmd(user, commentId, ctx.Param("resource_id"))
	if err1 != nil {
		commonctl.SendBadRequestParam(ctx, err1)

		return
	}

	if err = ctl.commentService.ReportComment(ctx.Request.Context(), cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, nil)
	}
}
