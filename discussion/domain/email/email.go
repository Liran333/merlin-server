package email

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	discussionprimitive "github.com/openmerlin/merlin-server/discussion/domain/primitive"
)

type ReportEmailParam struct {
	User          primitive.Account
	Index         domain.CodeRepoIndex
	Content       discussionprimitive.CommentContent
	IssueId       int64
	ReportType    string
	ResourceType  string
	ReportContent discussionprimitive.CommentContent
}

type Email interface {
	SendReportEmail(param ReportEmailParam) error
}
