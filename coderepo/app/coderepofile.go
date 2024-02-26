/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/coderepo/domain/repofileadapter"
)

// CodeRepoFileAppService is an interface for code repository file application service.
type CodeRepoFileAppService interface {
	List(cmd *CmdToFile) (*domain.ListFileInfo, error)
	Get(cmd *CmdToFile) (*domain.DetailFileInfo, error)
	Download(cmd *CmdToFile) (*domain.DownLoadFileInfo, error)
}

// NewCodeRepoFileAppService creates a new instance of CodeRepoFileAppService.
func NewCodeRepoFileAppService(repoFileAdapter repofileadapter.CodeRepoFileAdapter) *codeRepoFileAppService {
	return &codeRepoFileAppService{repoFileAdapter: repoFileAdapter}
}

type codeRepoFileAppService struct {
	repoFileAdapter repofileadapter.CodeRepoFileAdapter
}

// List lists the files in the code repository.
func (s *codeRepoFileAppService) List(cmd *CmdToFile) (*domain.ListFileInfo, error) {
	codeRepoFile := cmd.toCodeRepoFile()
	listFileInfo, err := s.repoFileAdapter.List(&codeRepoFile)
	if err != nil {
		return nil, err
	}
	return listFileInfo, nil
}

// Get gets the details of a file in the code repository.
func (s *codeRepoFileAppService) Get(cmd *CmdToFile) (*domain.DetailFileInfo, error) {
	codeRepoFile := cmd.toCodeRepoFile()
	detailFileInfo, err := s.repoFileAdapter.Get(&codeRepoFile)
	if err != nil {
		return nil, err
	}
	return detailFileInfo, nil
}

// Download downloads a file from the code repository.
func (s *codeRepoFileAppService) Download(cmd *CmdToFile) (*domain.DownLoadFileInfo, error) {
	codeRepoFile := cmd.toCodeRepoFile()
	contents, err := s.repoFileAdapter.Download(&codeRepoFile)
	if err != nil {
		return nil, err
	}
	return contents, nil
}
