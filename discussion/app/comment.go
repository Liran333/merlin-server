package app

import (
	"context"

	"golang.org/x/xerrors"

	"github.com/openmerlin/merlin-server/coderepo/domain/resourceadapter"
	"github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/discussion/domain"
	"github.com/openmerlin/merlin-server/discussion/domain/repository"
)

type CommentService interface {
	CreateIssueComment(context.Context, CmdToCreateIssueComment) error
	UpdateIssueComment(context.Context, CmdToUpdateIssueComment) error
	DeleteIssueComment(context.Context, CmdToDeleteIssueComment) error
	ReportComment(context.Context, CmdToReportComment) error
}

func NewCommentService(
	re resourceadapter.ResourceAdapter,
	p app.ResourcePermissionAppService,
	i repository.Issue,
	c repository.IssueComment,
	r repository.IssueCommentReport,
) *commentService {
	return &commentService{
		resource:          re,
		permission:        p,
		issueRepo:         i,
		commentRepo:       c,
		commentReportRepo: r,
	}
}

type commentService struct {
	resource          resourceadapter.ResourceAdapter
	permission        app.ResourcePermissionAppService
	issueRepo         repository.Issue
	commentRepo       repository.IssueComment
	commentReportRepo repository.IssueCommentReport
}

func (i *commentService) CreateIssueComment(ctx context.Context, cmd CmdToCreateIssueComment) error {
	r, err := i.resource.GetByIndex(cmd.ResourceId)
	if err != nil {
		return allerror.New(allerror.ErrorCodeRepoNotFound, "resource not found", err)
	}

	if err = i.permission.CanRead(ctx, cmd.Owner, r); err != nil {
		return err
	}

	issue, err := i.issueRepo.Find(cmd.IssueId)
	if err != nil {
		return allerror.NewNotFound(
			allerror.ErrorCodeIssueNotFound,
			"not found",
			xerrors.Errorf("failed to find issue by id, %w", err),
		)
	}

	if !issue.AllowComment() {
		return allerror.New(allerror.ErrorCodeIssueClosed, "issue is closed", err)
	}

	comment := domain.NewIssueComment(cmd.Owner, cmd.IssueId, cmd.Content)
	if err = i.commentRepo.Save(comment); err != nil {
		return allerror.New(allerror.ErrorCodeFailToCreateComment, "failed to create comment", err)
	}

	return nil
}

func (i *commentService) UpdateIssueComment(ctx context.Context, cmd CmdToUpdateIssueComment) error {
	r, err := i.resource.GetByIndex(cmd.ResourceId)
	if err != nil {
		return allerror.New(allerror.ErrorCodeRepoNotFound, "resource not found", err)
	}

	if err = i.permission.CanRead(ctx, cmd.User, r); err != nil {
		return err
	}

	//todo sensitive check
	comment, err := i.commentRepo.Find(cmd.CommentId)
	if err != nil {
		return allerror.NewNotFound(
			allerror.ErrorCodeCommentNotFound,
			"not found",
			xerrors.Errorf("failed to find comment by id, %w", err),
		)
	}

	if err = comment.UpdateContent(cmd.User, cmd.Content); err != nil {
		return err
	}

	if err = i.commentRepo.Save(comment); err != nil {
		return allerror.New(allerror.ErrorCodeFailToUpdateComment, "failed to update comment", err)
	}

	return nil
}

func (i *commentService) DeleteIssueComment(ctx context.Context, cmd CmdToDeleteIssueComment) error {
	r, err := i.resource.GetByIndex(cmd.ResourceId)
	if err != nil {
		return allerror.New(allerror.ErrorCodeRepoNotFound, "resource not found", err)
	}

	if err = i.permission.CanRead(ctx, cmd.User, r); err != nil {
		return err
	}

	comment, err := i.commentRepo.Find(cmd.CommentId)
	if err != nil {
		return allerror.NewNotFound(
			allerror.ErrorCodeCommentNotFound,
			"not found",
			xerrors.Errorf("failed to find comment by id, %w", err),
		)
	}

	if !comment.IsCommentOwner(cmd.User) && !i.isSecurity(cmd.User) {
		return allerror.NewNoPermission("no permission", xerrors.Errorf("cant delete"))
	}

	return i.commentRepo.Delete(comment.Id)
}

func (i *commentService) ReportComment(ctx context.Context, cmd CmdToReportComment) error {
	r, err := i.resource.GetByIndex(cmd.ResourceId)
	if err != nil {
		return allerror.New(allerror.ErrorCodeRepoNotFound, "resource not found", err)
	}

	if err = i.permission.CanRead(ctx, cmd.User, r); err != nil {
		return err
	}

	return i.commentReportRepo.Save(cmd.IssueCommentReport)
}

func (i *commentService) isSecurity(user primitive.Account) bool {
	//todo check security user
	return false
}
