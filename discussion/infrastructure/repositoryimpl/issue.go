package repositoryimpl

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/discussion/app"
	"github.com/openmerlin/merlin-server/discussion/domain"
	discussionprimitive "github.com/openmerlin/merlin-server/discussion/domain/primitive"
	"github.com/openmerlin/merlin-server/discussion/domain/repository"
)

const (
	fieldId     = "id"
	fieldStatus = "status"
)

func NewIssueImpl(db postgresql.Impl) *issueImpl {
	issueTableName = db.TableName()
	err := db.DB().AutoMigrate(&IssueDO{})
	if err != nil {
		logrus.Fatalf("failed to auto migrate %s table: %v", issueTableName, err)
	}

	return &issueImpl{Impl: db}
}

type issueImpl struct {
	postgresql.Impl
}

func (impl *issueImpl) Save(issue domain.Issue) (int64, error) {
	do := toIssueDO(issue)
	err := impl.DB().Save(&do).Error

	return do.Id, err
}

func (impl *issueImpl) Find(ctx context.Context, issueId int64) (issue domain.Issue, err error) {
	do := IssueDO{Id: issueId}
	if err = impl.GetByPrimaryKey(ctx, &do); err != nil {
		return
	}

	issue = do.toIssue()

	return
}

func (impl *issueImpl) List(resourceId primitive.Identity, option repository.IssueListOption,
) (data []app.IssueDTO, err error) {
	do := IssueDO{
		ResourceId: resourceId.Integer(),
	}

	if option.Status != nil {
		do.Status = option.Status.IssueStatus()
	}

	limit, offset := option.Paginate()

	var list []IssueDO
	err = impl.DB().Order(impl.OrderByDesc(fieldId)).Limit(limit).Offset(offset).Find(&list, &do).Error
	if err != nil {
		return
	}

	for _, v := range list {
		data = append(data, v.toIssueDTO())
	}

	return
}

type CountResult struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

func (impl *issueImpl) CountByStatus(resourceId primitive.Identity) (count app.ListIssuesCountDTO, err error) {
	do := IssueDO{
		ResourceId: resourceId.Integer(),
	}

	var results []CountResult
	err = impl.DB().Select("status, count(status) as count").Where(&do).Group(fieldStatus).Scan(&results).Error

	for _, v := range results {
		switch v.Status {
		case discussionprimitive.IssueStatusOpen.IssueStatus():
			count.Open = v.Count
		case discussionprimitive.IssueStatusClosed.IssueStatus():
			count.Closed = v.Count
		default:

		}
	}

	count.All = count.Open + count.Closed

	return
}
