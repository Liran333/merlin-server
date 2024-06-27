package repositoryimpl

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/discussion/domain"
)

func NewIssueCommentImpl(db postgresql.Impl) *issueCommentImpl {
	issueCommentTableName = db.TableName()
	err := db.DB().AutoMigrate(&IssueCommentDO{})
	if err != nil {
		logrus.Fatalf("failed to auto migrate %s table: %v", issueCommentTableName, err)
	}

	return &issueCommentImpl{Impl: db}
}

type issueCommentImpl struct {
	postgresql.Impl
}

func (impl *issueCommentImpl) Save(comment domain.IssueComment) (domain.IssueComment, error) {
	do := toIssueCommentDO(comment)

	err := impl.DB().Save(&do).Error

	return do.toIssueComment(), err
}

func (impl *issueCommentImpl) Find(ctx context.Context, commentId int64) (comment domain.IssueComment, err error) {
	do := IssueCommentDO{Id: commentId}
	if err = impl.GetByPrimaryKey(ctx, &do); err != nil {
		return
	}

	comment = do.toIssueComment()

	return
}

func (impl *issueCommentImpl) List(issueId int64) (comments []domain.IssueComment, err error) {
	do := IssueCommentDO{IssueId: issueId}

	var list []IssueCommentDO
	err = impl.DB().Find(&list, &do).Error
	if err != nil {
		return
	}

	for _, v := range list {
		comments = append(comments, v.toIssueComment())
	}

	return
}

func (impl *issueCommentImpl) Delete(ctx context.Context, commentId int64) error {
	do := IssueCommentDO{Id: commentId}
	return impl.DeleteByPrimaryKey(ctx, &do)
}
