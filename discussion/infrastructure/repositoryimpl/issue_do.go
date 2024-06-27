package repositoryimpl

import (
	"time"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/discussion/app"
	"github.com/openmerlin/merlin-server/discussion/domain"
	discussionprimitive "github.com/openmerlin/merlin-server/discussion/domain/primitive"
)

var issueTableName string

type IssueDO struct {
	Id           int64              `gorm:"primaryKey;autoIncrement"`
	Author       string             `gorm:"column:author"`
	Title        string             `gorm:"column:title"`
	Status       string             `gorm:"column:status"`
	Operation    []domain.Operation `gorm:"column:operation;serializer:json"`
	ResourceId   int64              `gorm:"column:resource_id;index"`
	ResourceType string             `gorm:"column:resource_type"`
	CommentCount int64              `gorm:"column:comment_count"`
	CreatedAt    time.Time          `gorm:"column:created_at;<-:create"`
	UpdatedAt    time.Time          `gorm:"column:updated_at;<-:update"`
}

func (do IssueDO) TableName() string {
	return issueTableName
}

func toIssueDO(issue domain.Issue) IssueDO {
	return IssueDO{
		Id:           issue.Id,
		Author:       issue.Author.Account(),
		Title:        issue.Title.Title(),
		Status:       issue.Status.IssueStatus(),
		ResourceId:   issue.Resource.Id.Integer(),
		ResourceType: string(issue.Resource.Type),
		Operation:    issue.Operation,
		CommentCount: issue.CommentCount,
	}
}

func (do IssueDO) toIssue() domain.Issue {
	return domain.Issue{
		Id:           do.Id,
		Author:       primitive.CreateAccount(do.Author),
		Title:        discussionprimitive.CreateIssueTitle(do.Title),
		Status:       discussionprimitive.CreateIssueStatus(do.Status),
		Operation:    do.Operation,
		CommentCount: do.CommentCount,
		CreatedAt:    do.CreatedAt,
		Resource: domain.Resource{
			Id:   primitive.CreateIdentity(do.ResourceId),
			Type: primitive.ObjType(do.ResourceType),
		},
	}
}

func (do IssueDO) toIssueDTO() app.IssueDTO {
	return app.IssueDTO{
		Id:           do.Id,
		Title:        do.Title,
		Owner:        do.Author,
		Status:       do.Status,
		CommentCount: do.CommentCount,
		CreatedAt:    do.CreatedAt.In(time.UTC).Format(app.TimeFormat),
	}
}
