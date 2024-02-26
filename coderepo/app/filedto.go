/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// CmdToFile is a struct representing a command to create or update a file.
type CmdToFile struct {
	Owner    primitive.Account
	Name     primitive.MSDName
	Ref      primitive.FileRef
	FilePath primitive.FilePath
}

func (cmd *CmdToFile) toCodeRepoFile() domain.CodeRepoFile {
	return domain.CodeRepoFile{
		Name:     cmd.Name,
		Owner:    cmd.Owner,
		Ref:      cmd.Ref,
		FilePath: cmd.FilePath,
	}
}
