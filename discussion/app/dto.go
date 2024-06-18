package app

import (
	"math"
	"time"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/discussion/domain"
	discussionprimitive "github.com/openmerlin/merlin-server/discussion/domain/primitive"
	"github.com/openmerlin/merlin-server/discussion/domain/repository"
)

const (
	timeFormat = "2006-01-02T15:04:05Z"
)

type CmdToCreateIssue struct {
	Resource domain.Resource
	Owner    primitive.Account
	Title    discussionprimitive.IssueTitle
	Content  discussionprimitive.CommentContent
}

type CmdToCloseIssue struct {
	User     primitive.Account
	Resource domain.Resource
	IssueId  int64
}

type CmdToReopenIssue = CmdToCloseIssue

type CmdToGetIssue struct {
	User       primitive.Account
	ResourceId primitive.Identity
	IssueId    int64

	PageNum      int
	CountPerPage int
}

type CmdToListIssues struct {
	ResourceId primitive.Identity
	Option     repository.IssueListOption
}

type ListIssuesCountDTO struct {
	All    int64 `json:"all"`
	Open   int64 `json:"open"`
	Closed int64 `json:"closed"`
}

type ListIssuesDTO struct {
	List []IssueDTO `json:"list"`
}

type IssueDTO struct {
	Id           int64  `json:"id"`
	Title        string `json:"title"`
	Owner        string `json:"owner"`
	Status       string `json:"status"`
	CommentCount int64  `json:"comment_count"`
	CreatedAt    string `json:"created_at"`
}

func ToIssueDTO(issue domain.Issue) IssueDTO {
	return IssueDTO{
		Id:           issue.Id,
		Title:        issue.Title.Title(),
		Owner:        issue.Author.Account(),
		Status:       issue.Status.IssueStatus(),
		CommentCount: issue.CommentCount,
		CreatedAt:    issue.CreatedAt.Format(timeFormat),
	}
}

type IssueDetailDTO struct {
	IsSecurity bool     `json:"is_security"`
	IsOwner    bool     `json:"is_owner"`
	Issue      IssueDTO `json:"issue"`
	Items      ItemsDTO `json:"items"`
}

type ItemsDTO []ItemDTO

func (d ItemsDTO) Len() int {
	return len(d)
}

func (d ItemsDTO) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d ItemsDTO) Less(i, j int) bool {
	return d[i].createdAt.Before(d[j].createdAt)
}

type ItemDTO struct {
	Id        int64  `json:"id"`
	Type      string `json:"type"`
	Owner     string `json:"owner"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	createdAt time.Time
}

func mergeOperationAndComments(operations []domain.Operation, comments []domain.IssueComment) ItemsDTO {
	var data ItemsDTO
	for _, v := range operations {
		data = append(data, ItemDTO{
			Type:      "operation",
			Owner:     v.User,
			Content:   v.Action,
			createdAt: v.CreatedAt,
			CreatedAt: v.CreatedAt.Format(timeFormat),
		})
	}

	for _, v := range comments {
		data = append(data, ItemDTO{
			Id:        v.Id,
			Type:      "comment",
			Owner:     v.Author.Account(),
			Content:   v.Content.CommentContent(),
			createdAt: v.CreatedAt,
			CreatedAt: v.CreatedAt.Format(timeFormat),
		})
	}

	return data
}

func (d ItemsDTO) paginate(pageNum, countPerPage int) ItemsDTO {
	start, end := d.slicePage(pageNum, countPerPage, len(d))

	return d[start:end]
}

func (d ItemsDTO) slicePage(page, pageSize, nums int) (sliceStart int, sliceEnd int) {
	if page < 0 {
		page = 1
	}

	if pageSize < 0 {
		pageSize = 10
	}

	if pageSize > nums {
		return 0, nums
	}

	pageCount := int(math.Ceil(float64(nums) / float64(pageSize)))
	if page > pageCount {
		return 0, 0
	}
	sliceStart = (page - 1) * pageSize
	sliceEnd = sliceStart + pageSize

	if sliceEnd > nums {
		sliceEnd = nums
	}

	return sliceStart, sliceEnd
}

type CmdToCreateIssueComment struct {
	IssueId    int64
	ResourceId primitive.Identity
	Owner      primitive.Account
	Content    discussionprimitive.CommentContent
}

type CmdToUpdateIssueComment struct {
	CommentId  int64
	ResourceId primitive.Identity
	Content    discussionprimitive.CommentContent
	User       primitive.Account
}

type CmdToDeleteIssueComment struct {
	CommentId  int64
	ResourceId primitive.Identity
	User       primitive.Account
}

type CmdToReportComment struct {
	domain.IssueCommentReport
	ResourceId primitive.Identity
}
