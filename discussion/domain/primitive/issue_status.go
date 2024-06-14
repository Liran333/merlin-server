package primitive

import "errors"

const (
	statusOpen   = "open"
	statusClosed = "closed"

	IssueStatusOpen   = issueStatus(statusOpen)
	IssueStatusClosed = issueStatus(statusClosed)
)

type IssueStatus interface {
	IssueStatus() string
	IsOpen() bool
}

func NewIssueStatus(v string) (IssueStatus, error) {
	if v != statusOpen && v != statusClosed {
		return nil, errors.New("invalid status")
	}

	return issueStatus(v), nil
}

func CreateIssueStatus(v string) IssueStatus {
	return issueStatus(v)
}

type issueStatus string

func (i issueStatus) IssueStatus() string {
	return string(i)
}

func (i issueStatus) IsOpen() bool {
	return i.IssueStatus() == statusOpen
}
