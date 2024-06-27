package repository

import (
	"context"

	"github.com/openmerlin/merlin-server/discussion/domain"
	"github.com/openmerlin/merlin-server/discussion/domain/primitive"
)

type IssueListOption struct {
	Status primitive.IssueStatus

	PageNum      int
	CountPerPage int
}

func (i IssueListOption) Paginate() (int, int) {
	offset := (i.PageNum - 1) * i.CountPerPage

	return i.CountPerPage, offset
}

type Issue interface {
	Save(issue domain.Issue) (int64, error)
	Find(context.Context, int64) (domain.Issue, error)
}

type IssueComment interface {
	Save(comment domain.IssueComment) (domain.IssueComment, error)
	List(int64) ([]domain.IssueComment, error)
	Delete(context.Context, int64) error
	Find(context.Context, int64) (domain.IssueComment, error)
}
