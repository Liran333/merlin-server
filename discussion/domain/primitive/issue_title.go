package primitive

import (
	"errors"

	"github.com/openmerlin/merlin-server/utils"
)

type IssueTitle interface {
	Title() string
}

func NewIssueTitle(v string) (IssueTitle, error) {
	if utils.StrLen(v) > cfg.MaxTitleLength ||
		utils.StrLen(utils.XSSEscapeString(v)) > cfg.MaxTitleLength {
		return nil, errors.New("title is too long")
	}

	return issueTitle(v), nil
}

func CreateIssueTitle(v string) IssueTitle {
	return issueTitle(v)
}

type issueTitle string

func (i issueTitle) Title() string {
	return string(i)
}
