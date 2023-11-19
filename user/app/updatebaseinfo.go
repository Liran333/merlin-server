package app

import "github.com/openmerlin/merlin-server/user/domain"

func (s userService) UpdateBasicInfo(account domain.Account, cmd UpdateUserBasicInfoCmd) error {
	user, err := s.repo.GetByAccount(account)
	if err != nil {
		return err
	}

	if b := cmd.toUser(&user); !b {
		return nil
	}

	if _, err = s.repo.Save(&user); err != nil {
		return err
	}

	return nil
}
