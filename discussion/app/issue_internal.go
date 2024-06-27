package app

import (
	"context"

	"golang.org/x/xerrors"

	"github.com/openmerlin/merlin-server/discussion/domain/repository"
)

type IssueInternalService interface {
	UpdateCommentCount(context.Context, int64, int64) error
}

func NewIssueInternalService(i repository.Issue, c repository.IssueComment) *issueInternalService {
	return &issueInternalService{issueRepo: i, comment: c}
}

type issueInternalService struct {
	issueRepo repository.Issue
	comment   repository.IssueComment
}

func (d *issueInternalService) UpdateCommentCount(ctx context.Context, issueId, increaseCount int64) error {
	issue, err := d.issueRepo.Find(ctx, issueId)
	if err != nil {
		return xerrors.Errorf("find issue %d failed: %w", issueId, err)
	}

	comments, err := d.comment.List(issueId)
	if err != nil {
		return xerrors.Errorf("find comments of %d failed: %w", issueId, err)
	}

	// the first comment does not count
	issue.SetCommentCount(int64(len(comments) - 1))
	if _, err = d.issueRepo.Save(issue); err != nil {
		return xerrors.Errorf("save issue %d failed: %w", issueId, err)
	}

	return nil
}
