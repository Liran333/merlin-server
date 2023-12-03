package app

import (
	"github.com/openmerlin/merlin-server/user/domain"
)

type UserInfoDTO struct {
	Points int `json:"points"`

	UserDTO
}

type UserDTO struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Account string `json:"account"`

	Bio      string `json:"bio"`
	AvatarId string `json:"avatar_id"`
}

func newUserDTO(u *domain.User) (dto UserDTO) {
	dto.Account = u.Account.Account()
	dto.AvatarId = u.AvatarId.AvatarId()
	dto.Bio = u.Bio.Bio()
	dto.Email = u.Email.Email()

	return
}

type UpdateUserBasicInfoCmd struct {
	Bio           domain.Bio
	Email         domain.Email
	AvatarId      domain.AvatarId
	bioChanged    bool
	avatarChanged bool
	emailChanged  bool
}

func (cmd *UpdateUserBasicInfoCmd) toUser(u *domain.User) (changed bool) {
	if cmd.AvatarId != nil && !domain.IsSameDomainValue(cmd.AvatarId, u.AvatarId) {
		u.AvatarId = cmd.AvatarId
		cmd.avatarChanged = true
	}

	if cmd.Bio != nil && !domain.IsSameDomainValue(cmd.Bio, u.Bio) {
		u.Bio = cmd.Bio
		cmd.bioChanged = true
	}

	if cmd.Email != nil && u.Email.Email() != cmd.Email.Email() {
		u.Email = cmd.Email
		cmd.emailChanged = true
	}

	changed = cmd.avatarChanged || cmd.bioChanged || cmd.emailChanged

	return
}

type SendBindEmailCmd struct {
	User  domain.Account
	Email domain.Email
	Capt  string
}
