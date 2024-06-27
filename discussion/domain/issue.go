package domain

import (
	"errors"
	"time"

	"golang.org/x/xerrors"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	discussionprimitive "github.com/openmerlin/merlin-server/discussion/domain/primitive"
)

const (
	operationReopen = "reopen"
	operationClose  = "close"
)

type Issue struct {
	Id           int64
	Title        discussionprimitive.IssueTitle
	Author       primitive.Account
	Status       discussionprimitive.IssueStatus
	Operation    []Operation
	Resource     Resource
	CommentCount int64
	CreatedAt    time.Time
}

type Resource struct {
	Id   primitive.Identity
	Type primitive.ObjType
}

type Operation struct {
	User      string
	Action    string
	CreatedAt time.Time
}

func NewIssue(
	resource Resource,
	author primitive.Account,
	title discussionprimitive.IssueTitle,
) Issue {
	return Issue{
		Title:    title,
		Author:   author,
		Status:   discussionprimitive.IssueStatusOpen,
		Resource: resource,
	}
}

func (i *Issue) Close(user primitive.Account) error {
	if !i.Status.IsOpen() {
		return allerror.New(
			allerror.ErrorCodeIssueClosed,
			"failed to close issue",
			errors.New("issue is closed"),
		)
	}

	i.Status = discussionprimitive.IssueStatusClosed

	i.Operation = append(i.Operation, Operation{
		User:      user.Account(),
		Action:    operationClose,
		CreatedAt: time.Now(),
	})

	return nil
}

func (i *Issue) Reopen(user primitive.Account) error {
	if i.Status.IsOpen() {
		return allerror.New(
			allerror.ErrorCodeIssueIsOpen,
			"failed to reopen issue",
			errors.New("issue is open"),
		)
	}

	i.Status = discussionprimitive.IssueStatusOpen

	i.Operation = append(i.Operation, Operation{
		User:      user.Account(),
		Action:    operationReopen,
		CreatedAt: time.Now(),
	})

	return nil
}

func (i *Issue) IsStatusChanged(status discussionprimitive.IssueStatus) bool {
	return status != i.Status
}

func (i *Issue) SetCommentCount(count int64) {
	if count < 0 {
		count = 0
	}

	i.CommentCount = count
}

func (i *Issue) AllowComment() error {
	if i.Status.IsOpen() {
		return nil
	}

	return allerror.New(
		allerror.ErrorCodeIssueClosed,
		"issue is closed",
		xerrors.Errorf("issue is closed, cant comment"),
	)
}

func (i *Issue) IsIssueAuthor(user primitive.Account) bool {
	return i.Author == user
}
