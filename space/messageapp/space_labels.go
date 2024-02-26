/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package messageapp

import (
	"github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/space/domain/repository"
)

// SpaceAppService is an interface for space application service.
type SpaceAppService interface {
	ResetLabels(*domain.SpaceIndex, *CmdToResetLabels) error
}

// NewSpaceAppService creates a new instance of SpaceAppService.
func NewSpaceAppService(repoAdapter repository.SpaceLabelsRepoAdapter) SpaceAppService {
	return &spaceAppService{
		repoAdapter: repoAdapter,
	}
}

type spaceAppService struct {
	repoAdapter repository.SpaceLabelsRepoAdapter
}

// ResetLabels resets labels for the given space index using the provided command.
func (s *spaceAppService) ResetLabels(index *domain.SpaceIndex, cmd *CmdToResetLabels) error {
	return s.repoAdapter.Save(index, cmd)
}
