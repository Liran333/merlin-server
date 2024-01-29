package app

import (
	"github.com/openmerlin/merlin-server/coderepo/domain/repository"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
)

type Permission interface {
	Check(primitive.Account, primitive.Account, primitive.ObjType, primitive.Action) error
}

type BranchAppService interface {
	Create(primitive.Account, *CmdToCreateBranch) (BranchCreateDTO, error)
	Delete(primitive.Account, *CmdToDeleteBranch) error
}

func NewBranchAppService(
	permission Permission,
	branchAdapter repository.BranchRepositoryAdapter,
	branchClientAdapter repository.BranchClientAdapter,
	checkRepoAdapter repository.CheckRepoAdapter,
) BranchAppService {
	return &branchAppService{
		permission:          permission,
		branchAdapter:       branchAdapter,
		branchClientAdapter: branchClientAdapter,
		checkRepoAdapter:    checkRepoAdapter,
	}
}

type branchAppService struct {
	permission          Permission
	branchAdapter       repository.BranchRepositoryAdapter
	branchClientAdapter repository.BranchClientAdapter
	checkRepoAdapter    repository.CheckRepoAdapter
}

func (s *branchAppService) Create(user primitive.Account, cmd *CmdToCreateBranch) (dto BranchCreateDTO, err error) {
	if user != cmd.Owner {
		err = s.permission.Check(
			user, cmd.Owner, primitive.ObjTypeModel, primitive.ActionCreate,
		)
		if err != nil {
			return
		}
	}

	if err = s.checkRepoAdapter.CheckRepo(cmd.RepoType, cmd.Owner, cmd.Repo); err != nil {
		return
	}

	branch := cmd.toBranch()

	v, err := s.branchClientAdapter.CreateBranch(&branch)
	if err != nil {
		return
	}

	if err = s.branchAdapter.Add(&branch); err == nil {
		dto = toBranchCreateDTO(v)
	}

	return
}

func (s *branchAppService) Delete(user primitive.Account, cmd *CmdToDeleteBranch) (err error) {
	if user != cmd.Owner {
		err = s.permission.Check(
			user, cmd.Owner, primitive.ObjTypeModel, primitive.ActionDelete,
		)
		if err != nil {
			return
		}
	}

	br, err := s.branchAdapter.FindByIndex(&cmd.BranchIndex)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}

	if err = s.branchClientAdapter.DeleteBranch(&cmd.BranchIndex); err != nil {
		return
	}

	err = s.branchAdapter.Delete(br.Id)

	return
}
