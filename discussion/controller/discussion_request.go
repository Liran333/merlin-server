package controller

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/discussion/app"
	"github.com/openmerlin/merlin-server/discussion/domain"
	discussionprimitive "github.com/openmerlin/merlin-server/discussion/domain/primitive"
	"github.com/openmerlin/merlin-server/discussion/domain/repository"
)

type reqToCreateIssue struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (r reqToCreateIssue) action() string {
	return fmt.Sprintf("create issue of %s/%s", r.Title, r.Content)
}

func (r reqToCreateIssue) toCreateIssueCmd(resourceId string, owner primitive.Account,
) (cmd app.CmdToCreateIssue, err error) {
	id, err := primitive.NewIdentity(resourceId)
	if err != nil {
		return
	}

	title, err := discussionprimitive.NewIssueTitle(r.Title)
	if err != nil {
		return
	}

	content, err := discussionprimitive.NewCommentContent(r.Content)
	if err != nil {
		return
	}

	cmd = app.CmdToCreateIssue{
		Resource: domain.Resource{
			Id: id,
		},
		Owner:   owner,
		Title:   title,
		Content: content,
	}

	return
}

func toCloseIssueCmd(user primitive.Account, resourceId string, issueId int64,
) (cmd app.CmdToCloseIssue, err error) {
	id, err := primitive.NewIdentity(resourceId)
	if err != nil {
		return
	}

	cmd = app.CmdToCloseIssue{
		User:    user,
		IssueId: issueId,
		Resource: domain.Resource{
			Id: id,
		},
	}

	return
}

type reqToListIssue struct {
	controller.CommonListRequest
	Status string `form:"status"`
}

func (r reqToListIssue) toListIssuesCmd(resourceId string) (cmd app.CmdToListIssues, err error) {
	id, err := primitive.NewIdentity(resourceId)
	if err != nil {
		return
	}

	status, _ := discussionprimitive.NewIssueStatus(r.Status)

	if r.PageNum <= 0 {
		r.PageNum = 1
	}

	if r.CountPerPage <= 0 {
		r.CountPerPage = 50
	}

	return app.CmdToListIssues{
		Resource: domain.Resource{
			Id: id,
		},
		Option: repository.IssueListOption{
			Status:       status,
			PageNum:      r.PageNum,
			CountPerPage: r.CountPerPage,
		},
	}, nil
}

type reqToGetIssue struct {
	controller.CommonListRequest
}

func (r reqToGetIssue) toGetIssueCmd(user primitive.Account, resourceId string, issueId int64,
) (cmd app.CmdToGetIssue, err error) {
	if r.PageNum <= 0 {
		r.PageNum = 1
	}

	if r.CountPerPage <= 0 {
		r.CountPerPage = 50
	}

	id, err := primitive.NewIdentity(resourceId)
	if err != nil {
		return
	}

	cmd = app.CmdToGetIssue{
		User:         user,
		IssueId:      issueId,
		PageNum:      r.PageNum,
		CountPerPage: r.CountPerPage,
		Resource: domain.Resource{
			Id: id,
		},
	}

	return
}

type reqToCreateComment struct {
	IssueId int64  `json:"issue_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (r reqToCreateComment) toCreateCommentCmd(user primitive.Account, resourceId string,
) (cmd app.CmdToCreateIssueComment, err error) {
	id, err := primitive.NewIdentity(resourceId)
	if err != nil {
		return
	}

	content, err := discussionprimitive.NewCommentContent(r.Content)
	if err != nil {
		return
	}

	cmd = app.CmdToCreateIssueComment{
		IssueId: r.IssueId,
		Owner:   user,
		Content: content,
		Resource: domain.Resource{
			Id: id,
		},
	}

	return
}

type reqToUpdateComment struct {
	Content string `json:"content" binding:"required"`
}

func (r reqToUpdateComment) toUpdateCommentCmd(commentId int64, user primitive.Account, resourceId string,
) (app.CmdToUpdateIssueComment, error) {
	id, err := primitive.NewIdentity(resourceId)
	if err != nil {
		return app.CmdToUpdateIssueComment{}, err
	}

	content, err := discussionprimitive.NewCommentContent(r.Content)
	if err != nil {
		return app.CmdToUpdateIssueComment{}, err
	}

	return app.CmdToUpdateIssueComment{
		CommentId: commentId,
		User:      user,
		Content:   content,
		Resource: domain.Resource{
			Id: id,
		},
	}, nil
}

func toDeleteCommentCmd(commentId int64, user primitive.Account, resourceId string,
) (app.CmdToDeleteIssueComment, error) {
	id, err := primitive.NewIdentity(resourceId)
	if err != nil {
		return app.CmdToDeleteIssueComment{}, err
	}

	return app.CmdToDeleteIssueComment{
		CommentId: commentId,
		User:      user,
		Resource: domain.Resource{
			Id: id,
		},
	}, nil
}

type reqToReportComment struct {
	Type    string `json:"type" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (r reqToReportComment) toReportCommentCmd(user primitive.Account, commentId int64, resourceId string,
) (app.CmdToReportComment, error) {
	id, err := primitive.NewIdentity(resourceId)
	if err != nil {
		return app.CmdToReportComment{}, err
	}

	content, err := discussionprimitive.NewCommentContent(r.Content)
	if err != nil {
		return app.CmdToReportComment{}, err
	}

	return app.CmdToReportComment{
		User:      user,
		Type:      r.Type,
		Content:   content,
		CommentId: commentId,
		Resource: domain.Resource{
			Id: id,
		},
	}, nil
}

type reqToUpdateCommentCount struct {
	Count int64 `json:"count" binding:"required"`
}
