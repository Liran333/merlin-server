package app

import (
	"context"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/openmerlin/merlin-server/coderepo/domain/resourceadapter"
	"github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/discussion/domain"
	"github.com/openmerlin/merlin-server/discussion/domain/email"
	"github.com/openmerlin/merlin-server/discussion/domain/message"
	"github.com/openmerlin/merlin-server/discussion/domain/repository"
)

type CommentService interface {
	CreateIssueComment(context.Context, CmdToCreateIssueComment) (ItemDTO, error)
	UpdateIssueComment(context.Context, CmdToUpdateIssueComment) error
	DeleteIssueComment(context.Context, CmdToDeleteIssueComment) error
	ReportComment(context.Context, CmdToReportComment) error
}

func NewCommentService(
	re resourceadapter.ResourceAdapter,
	p app.ResourcePermissionAppService,
	i repository.Issue,
	c repository.IssueComment,
	m message.CommentMessage,
	e email.Email,
) *commentService {
	rp := resourcePermission{
		resource:   re,
		permission: p,
	}

	return &commentService{
		issueRepo:          i,
		commentRepo:        c,
		message:            m,
		email:              e,
		resourcePermission: rp,
	}
}

type commentService struct {
	email              email.Email
	message            message.CommentMessage
	issueRepo          repository.Issue
	commentRepo        repository.IssueComment
	resourcePermission resourcePermission
}

func (i *commentService) CreateIssueComment(ctx context.Context, cmd CmdToCreateIssueComment) (ItemDTO, error) {
	if _, err := i.resourcePermission.CanRead(ctx, cmd.Resource.Id, cmd.Owner); err != nil {
		return ItemDTO{}, err
	}

	issue, err := i.issueRepo.Find(ctx, cmd.IssueId)
	if err != nil {
		return ItemDTO{}, allerror.NewNotFound(
			allerror.ErrorCodeIssueNotFound,
			"not found",
			xerrors.Errorf("failed to find issue by id, %w", err),
		)
	}

	if err = issue.AllowComment(); err != nil {
		return ItemDTO{}, err
	}

	comment := domain.NewIssueComment(cmd.Owner, cmd.IssueId, cmd.Content)
	savedComment, err := i.commentRepo.Save(comment)
	if err != nil {
		return ItemDTO{}, allerror.New(allerror.ErrorCodeFailToCreateComment, "failed to create comment", err)
	}

	event := domain.NewUpdateCommentCountEvent(cmd.IssueId, 1)
	if err = i.message.SendUpdateCommentCountEvent(event); err != nil {
		logrus.Errorf("send update comment count +1 of issue %d failed: %s", cmd.IssueId, err.Error())
	}

	return commentToItemDTO(savedComment), nil
}

func (i *commentService) UpdateIssueComment(ctx context.Context, cmd CmdToUpdateIssueComment) error {
	if _, err := i.resourcePermission.CanRead(ctx, cmd.Resource.Id, cmd.User); err != nil {
		return err
	}

	//todo sensitive check
	comment, err := i.commentRepo.Find(ctx, cmd.CommentId)
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

	if _, err = i.commentRepo.Save(comment); err != nil {
		return allerror.New(allerror.ErrorCodeFailToUpdateComment, "failed to update comment", err)
	}

	return nil
}

func (i *commentService) DeleteIssueComment(ctx context.Context, cmd CmdToDeleteIssueComment) error {
	if _, err := i.resourcePermission.CanRead(ctx, cmd.Resource.Id, cmd.User); err != nil {
		return err
	}

	comment, err := i.commentRepo.Find(ctx, cmd.CommentId)
	if err != nil {
		return allerror.NewNotFound(
			allerror.ErrorCodeCommentNotFound,
			"not found",
			xerrors.Errorf("failed to find comment by id, %w", err),
		)
	}

	if comment.IsFirstCommentOfIssue() {
		return allerror.New(
			allerror.ErrorCodeFailToDeleteComment,
			"cant delete first comment",
			xerrors.Errorf("cant delete first comment"),
		)
	}

	if !comment.IsCommentOwner(cmd.User) && !i.isSecurity(cmd.User) {
		return allerror.NewNoPermission("no permission", xerrors.Errorf("cant delete"))
	}

	if err = i.commentRepo.Delete(ctx, comment.Id); err != nil {
		return allerror.New(allerror.ErrorCodeFailToDeleteComment, "failed to delete comment", err)
	}

	event := domain.NewUpdateCommentCountEvent(comment.IssueId, -1)
	if err = i.message.SendUpdateCommentCountEvent(event); err != nil {
		logrus.Errorf("send update comment count -1 of issue %d failed: %s", comment.IssueId, err.Error())
	}

	return nil
}

func (i *commentService) ReportComment(ctx context.Context, cmd CmdToReportComment) error {
	r, err := i.resourcePermission.CanRead(ctx, cmd.Resource.Id, cmd.User)
	if err != nil {
		return err
	}

	comment, err := i.commentRepo.Find(ctx, cmd.CommentId)
	if err != nil {
		return allerror.NewNotFound(
			allerror.ErrorCodeCommentNotFound,
			"not found",
			xerrors.Errorf("failed to find comment by id, %w", err),
		)
	}

	issue, err := i.issueRepo.Find(ctx, comment.IssueId)
	if err != nil {
		return allerror.NewNotFound(
			allerror.ErrorCodeIssueNotFound,
			"not found",
			xerrors.Errorf("failed to find issue by id, %w", err),
		)
	}

	param := email.ReportEmailParam{
		User:          cmd.User,
		Index:         r.RepoIndex(),
		Content:       comment.Content,
		IssueId:       issue.Id,
		ReportType:    cmd.Type,
		ResourceType:  string(r.ResourceType()),
		ReportContent: cmd.Content,
	}

	return i.email.SendReportEmail(param)
}

func (i *commentService) isSecurity(user primitive.Account) bool {
	//todo check security user
	return false
}
