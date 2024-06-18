package app

import (
	"fmt"

	"github.com/openmerlin/merlin-server/discussion/domain/repository"
)

type IssueInternalService interface {
	UpdateCommentCount(int64, int64) error
}

func NewIssueInternalService(i repository.Issue) *issueInternalService {
	return &issueInternalService{issueRepo: i}
}

type issueInternalService struct {
	issueRepo repository.Issue
}

func (d *issueInternalService) UpdateCommentCount(issueId, increaseCount int64) error {
	issue, err := d.issueRepo.Find(issueId)
	if err != nil {
		return fmt.Errorf("find issue %d failed: %w", issueId, err)
	}

	issue.IncreaseCommentCount(increaseCount)

	_, err = d.issueRepo.Save(issue)

	return fmt.Errorf("save issue %d failed: %w", issueId, err)
}
