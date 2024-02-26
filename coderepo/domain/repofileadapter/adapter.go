/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repofileadapter provides an interface for interacting with code repository files.
package repofileadapter

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
)

// CodeRepoFileAdapter is an interface for interacting with code repository files.
type CodeRepoFileAdapter interface {
	List(*domain.CodeRepoFile) (*domain.ListFileInfo, error)

	Get(*domain.CodeRepoFile) (*domain.DetailFileInfo, error)

	Download(*domain.CodeRepoFile) (*domain.DownLoadFileInfo, error)
}
