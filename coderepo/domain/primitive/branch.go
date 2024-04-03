/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides primitive types and utility functions for working with basic concepts.
package primitive

import (
	"errors"
	"strings"
)

const (
	repoTypeModel = "model"
	repoTypeSpace = "space"
)

// Identity is an interface that represents an identity with both a string and integer representation.
type Identity interface {
	Identity() string
	Integer() int64
}

// BranchName represents a branch name.
type BranchName interface {
	BranchName() string
}

// NewBranchName creates a new BranchName from the given string value.
func NewBranchName(v string) (BranchName, error) {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "" {
		return nil, errors.New("branch name empty")
	}

	if len(v) > branchConfig.BranchNameMaxLength || len(v) < branchConfig.BranchNameMinLength {
		return nil, errors.New("branch name length is invalid")
	}

	if !branchConfig.branchRegexp.MatchString(v) {
		return nil, errors.New("branch name can only contain alphabet, integer, _ and -")
	}

	return branchName(v), nil
}

// CreateBranchName creates a new BranchName without validating the value.
func CreateBranchName(v string) BranchName {
	return branchName(v)
}

type branchName string

// BranchName returns the branch name as a string.
func (r branchName) BranchName() string {
	return string(r)
}

// RepoType represents a repository type.
type RepoType interface {
	RepoType() string
	IsModel() bool
	IsSpace() bool
}

// NewRepoType creates a new RepoType from the given string value.
func NewRepoType(v string) (RepoType, error) {
	v = strings.ToLower(strings.TrimSpace(v))

	if v != repoTypeModel && v != repoTypeSpace {
		return nil, errors.New("repo type incorrect")
	}

	return repoType(v), nil
}

// CreateRepoType creates a new RepoType without validating the value.
func CreateRepoType(v string) RepoType {
	return repoType(v)
}

type repoType string

// RepoType returns the repo type as a string.
func (r repoType) RepoType() string {
	return string(r)
}

// IsModel checks if the repo type is a model.
func (r repoType) IsModel() bool {
	return string(r) == repoTypeModel
}

// IsSpace checks if the repo type is a space.
func (r repoType) IsSpace() bool {
	return string(r) == repoTypeSpace
}
