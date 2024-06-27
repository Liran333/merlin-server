package repositoryimpl

type Tables struct {
	Issue        string `json:"issue" required:"true"`
	IssueComment string `json:"issue_comment" required:"true"`
}
