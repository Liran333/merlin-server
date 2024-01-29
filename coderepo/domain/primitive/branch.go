package primitive

import (
	"errors"
	"regexp"
	"strings"
)

const (
	repoTypeModel = "model"
	repoTypeSpace = "space"
)

var (
	regBranchName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// branch name
type BranchName interface {
	BranchName() string
}

func NewBranchName(v string) (BranchName, error) {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "" {
		return nil, errors.New("branch name empty")
	}

	if len(v) > maxBranchNameLength {
		return nil, errors.New("branch name too long")
	}

	if !regBranchName.MatchString(v) {
		return nil, errors.New("branch name can only contain alphabet, integer, _ and -")
	}

	return branchName(v), nil
}

func CreateBranchName(v string) BranchName {
	return branchName(v)
}

type branchName string

func (r branchName) BranchName() string {
	return string(r)
}

// repo type
type RepoType interface {
	RepoType() string
	IsModel() bool
	IsSpace() bool
}

func NewRepoType(v string) (RepoType, error) {
	v = strings.ToLower(strings.TrimSpace(v))

	if v != repoTypeModel && v != repoTypeSpace {
		return nil, errors.New("repo type incorrect")
	}

	return repoType(v), nil
}

func CreateRepoType(v string) RepoType {
	return repoType(v)
}

type repoType string

func (r repoType) RepoType() string {
	return string(r)
}

func (r repoType) IsModel() bool {
	return string(r) == repoTypeModel
}

func (r repoType) IsSpace() bool {
	return string(r) == repoTypeSpace
}
