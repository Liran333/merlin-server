package domain

import (
	"time"

	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	discussionprimitive "github.com/openmerlin/merlin-server/discussion/domain/primitive"
)

const (
	operationReopen = "reopen"
	operationClose  = "close"
)

type Index = domain.CodeRepoIndex

type Issue struct {
	Id           int64
	Index        Index
	Title        discussionprimitive.IssueTitle
	Author       primitive.Account
	Status       discussionprimitive.IssueStatus
	Operation    []Operation
	CommentCount int64
	CreatedAt    time.Time
}

type Operation struct {
	User      string
	Action    string
	CreatedAt time.Time
}

func NewIssue(index Index, author primitive.Account, title discussionprimitive.IssueTitle) Issue {
	return Issue{
		Index:  index,
		Title:  title,
		Author: author,
		Status: discussionprimitive.IssueStatusOpen,
	}
}

func (i *Issue) Close(user primitive.Account) {
	i.Status = discussionprimitive.IssueStatusClosed

	i.Operation = append(i.Operation, Operation{
		User:      user.Account(),
		Action:    operationClose,
		CreatedAt: time.Now(),
	})
}

func (i *Issue) Reopen(user primitive.Account) {
	i.Status = discussionprimitive.IssueStatusOpen

	i.Operation = append(i.Operation, Operation{
		User:      user.Account(),
		Action:    operationReopen,
		CreatedAt: time.Now(),
	})
}

func (i *Issue) IsStatusChanged(status discussionprimitive.IssueStatus) bool {
	return status != i.Status
}

func (i *Issue) IncreaseCommentCount(count int64) {
	i.CommentCount = i.CommentCount + count
}

func (i *Issue) AllowComment() bool {
	return i.Status.IsOpen()
}

func (i *Issue) IsIssueAuthor(user primitive.Account) bool {
	return i.Author == user
}
