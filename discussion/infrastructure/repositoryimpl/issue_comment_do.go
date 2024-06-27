package repositoryimpl

import (
	"time"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/discussion/domain"
	discussionprimitive "github.com/openmerlin/merlin-server/discussion/domain/primitive"
)

var issueCommentTableName string

type IssueCommentDO struct {
	Id int64 `gorm:"primaryKey;autoIncrement"`

	Author         string    `gorm:"column:author"`
	IssueId        int64     `gorm:"column:issue_id;index"`
	Content        string    `gorm:"column:content"`
	IsFirstComment bool      `gorm:"column:is_first_comment"`
	CreatedAt      time.Time `gorm:"column:created_at;<-:create"`
	UpdatedAt      time.Time `gorm:"column:updated_at;<-:update"`
}

func (do IssueCommentDO) TableName() string {
	return issueCommentTableName
}

func toIssueCommentDO(comment domain.IssueComment) IssueCommentDO {
	return IssueCommentDO{
		Id:             comment.Id,
		Author:         comment.Author.Account(),
		IssueId:        comment.IssueId,
		Content:        comment.Content.CommentContent(),
		IsFirstComment: comment.IsFirstComment,
	}
}

func (do IssueCommentDO) toIssueComment() domain.IssueComment {
	return domain.IssueComment{
		Id:             do.Id,
		Author:         primitive.CreateAccount(do.Author),
		IssueId:        do.IssueId,
		Content:        discussionprimitive.CreateCommentContent(do.Content),
		CreatedAt:      do.CreatedAt,
		IsFirstComment: do.IsFirstComment,
	}
}
