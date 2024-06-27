package app

import (
	"context"
	"errors"
	"sort"

	"golang.org/x/xerrors"

	"github.com/openmerlin/merlin-server/coderepo/domain/resourceadapter"
	"github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/discussion/domain"
	"github.com/openmerlin/merlin-server/discussion/domain/repository"
)

var discussionDisabledErr = allerror.New(
	allerror.ErrorCodeDiscussionDisabled, "discussion disabled", errors.New("discussion disabled"))

type IssueService interface {
	ListIssuesCount(context.Context, primitive.Account, primitive.Identity) (ListIssuesCountDTO, error)
	ListIssues(context.Context, primitive.Account, CmdToListIssues) (ListIssuesDTO, error)
	CreateIssue(context.Context, CmdToCreateIssue) error
	CloseIssue(context.Context, CmdToCloseIssue) error
	ReopenIssue(context.Context, CmdToReopenIssue) error
	GetIssue(context.Context, CmdToGetIssue) (IssueDetailDTO, error)
}

type IssueRepoQuery interface {
	List(primitive.Identity, repository.IssueListOption) ([]IssueDTO, error)
	CountByStatus(primitive.Identity) (count ListIssuesCountDTO, err error)
}

func NewIssueService(
	re resourceadapter.ResourceAdapter,
	p app.ResourcePermissionAppService,
	i repository.Issue,
	iq IssueRepoQuery,
	c repository.IssueComment,
) *issueService {
	rp := resourcePermission{
		resource:   re,
		permission: p,
	}

	return &issueService{
		resourcePermission: rp,
		issueRepo:          i,
		issueRepoQuery:     iq,
		commentRepo:        c,
	}
}

type issueService struct {
	resourcePermission resourcePermission
	issueRepo          repository.Issue
	issueRepoQuery     IssueRepoQuery
	commentRepo        repository.IssueComment
}

func (i *issueService) CreateIssue(ctx context.Context, cmd CmdToCreateIssue) error {
	//todo sensitive check

	_, err := i.resourcePermission.CanRead(ctx, cmd.Resource.Id, cmd.Owner)
	if err != nil {
		return err
	}

	issue := domain.NewIssue(cmd.Resource, cmd.Owner, cmd.Title)
	issueId, err := i.issueRepo.Save(issue)
	if err != nil {
		return allerror.New(allerror.ErrorCodeFailToCreateIssue, "failed to create issue", err)
	}

	comment := domain.NewFirstIssueComment(cmd.Owner, issueId, cmd.Content)
	if _, err = i.commentRepo.Save(comment); err != nil {
		return allerror.New(allerror.ErrorCodeFailToCreateComment, "failed to create comment", err)
	}

	return nil
}

func (i *issueService) CloseIssue(ctx context.Context, cmd CmdToCloseIssue) error {
	_, err := i.resourcePermission.CanRead(ctx, cmd.Resource.Id, cmd.User)
	if err != nil {
		return err
	}

	canUpdateResource := false
	if err = i.resourcePermission.CanUpdate(ctx, cmd.Resource.Id, cmd.User); err == nil {
		canUpdateResource = true
	}

	issue, err := i.issueRepo.Find(ctx, cmd.IssueId)
	if err != nil {
		return allerror.NewNotFound(
			allerror.ErrorCodeIssueNotFound,
			"not found",
			xerrors.Errorf("failed to find issue by id, %w", err),
		)
	}

	if !canUpdateResource && !issue.IsIssueAuthor(cmd.User) {
		return allerror.NewNoPermission("no permission", xerrors.Errorf("cant update"))
	}

	if err = issue.Close(cmd.User); err != nil {
		return err
	}

	if _, err = i.issueRepo.Save(issue); err != nil {
		return allerror.New(allerror.ErrorCodeFailToUpdateIssue, "failed to close issue", err)
	}

	return nil
}

func (i *issueService) ReopenIssue(ctx context.Context, cmd CmdToReopenIssue) error {
	_, err := i.resourcePermission.CanRead(ctx, cmd.Resource.Id, cmd.User)
	if err != nil {
		return err
	}

	canUpdateResource := false
	if err = i.resourcePermission.CanUpdate(ctx, cmd.Resource.Id, cmd.User); err == nil {
		canUpdateResource = true
	}

	issue, err := i.issueRepo.Find(ctx, cmd.IssueId)
	if err != nil {
		return allerror.NewNotFound(
			allerror.ErrorCodeIssueNotFound,
			"not found",
			xerrors.Errorf("failed to find issue by id, %w", err),
		)
	}

	if !canUpdateResource && !issue.IsIssueAuthor(cmd.User) {
		return allerror.NewNoPermission("no permission", xerrors.Errorf("cant update"))
	}

	if err = issue.Reopen(cmd.User); err != nil {
		return err
	}

	if _, err = i.issueRepo.Save(issue); err != nil {
		return allerror.New(allerror.ErrorCodeFailToUpdateIssue, "failed to update issue", err)
	}

	return nil
}

func (i *issueService) ListIssuesCount(ctx context.Context, user primitive.Account, id primitive.Identity,
) (dto ListIssuesCountDTO, err error) {
	_, err = i.resourcePermission.CanRead(ctx, id, user)
	if err != nil {
		return
	}

	return i.issueRepoQuery.CountByStatus(id)
}

func (i *issueService) ListIssues(ctx context.Context, user primitive.Account, cmd CmdToListIssues,
) (dto ListIssuesDTO, err error) {
	_, err = i.resourcePermission.CanRead(ctx, cmd.Resource.Id, user)
	if err != nil {
		return
	}

	issuesDTO, err := i.issueRepoQuery.List(cmd.Resource.Id, cmd.Option)

	return ListIssuesDTO{
		List: issuesDTO,
	}, err
}

func (i *issueService) GetIssue(ctx context.Context, cmd CmdToGetIssue) (dto IssueDetailDTO, err error) {
	_, err = i.resourcePermission.CanRead(ctx, cmd.Resource.Id, cmd.User)
	if err != nil {
		return
	}

	issue, err := i.issueRepo.Find(ctx, cmd.IssueId)
	if err != nil {
		return IssueDetailDTO{}, allerror.NewNotFound(
			allerror.ErrorCodeIssueNotFound,
			"not found",
			xerrors.Errorf("failed to find issue by id, %w", err),
		)
	}

	// check whether user has the permission to close/reopen issue
	var isOwner = false
	if err = i.resourcePermission.CanUpdate(ctx, cmd.Resource.Id, cmd.User); err == nil {
		isOwner = true
	}

	comments, err := i.commentRepo.List(cmd.IssueId)
	if err != nil {
		return IssueDetailDTO{}, xerrors.Errorf("find comments error: %w", err)
	}

	itemsDTO := mergeOperationAndComments(issue.Operation, comments)
	sort.Sort(itemsDTO)

	itemsDTOPaginate := itemsDTO.paginate(cmd.PageNum, cmd.CountPerPage)
	return IssueDetailDTO{
		IsSecurity: i.isSecurity(cmd.User),
		IsOwner:    isOwner || issue.IsIssueAuthor(cmd.User),
		Issue:      ToIssueDTO(issue),
		Items:      itemsDTOPaginate,
	}, nil
}

func (i *issueService) isSecurity(user primitive.Account) bool {
	//todo check security user
	return false
}
