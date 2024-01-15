package app

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
)

type CmdToCreateUser = domain.UserCreateCmd

type UserInfoDTO struct {
	UserDTO
}

type UserDTO struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Account string `json:"account"`

	Bio       string `json:"description"`
	Fullname  string `json:"fullname"`
	AvatarId  string `json:"avatar_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	Password  string `json:"-"`
}

type TokenDTO struct {
	Id         string `json:"id"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
	Expire     int64  `json:"expired"`
	Account    string `json:"account"`
	OwnerId    string `json:"owner_id"`
	Name       string `json:"name"`
	Permission string `json:"permission"`
	Token      string `json:"token"`
}

func newTokenDTO(t *domain.PlatformToken) (dto TokenDTO) {
	dto.CreatedAt = t.CreatedAt
	dto.UpdatedAt = t.UpdatedAt
	dto.Expire = t.Expire
	dto.Account = t.Account.Account()
	dto.Name = t.Name.Account()
	dto.Permission = t.Permission.TokenPerm()
	dto.Token = t.Token
	dto.Id = t.Id.Identity()
	dto.OwnerId = t.OwnerId.Identity()

	return

}

type AvatarDTO struct {
	AvatarId string `json:"avatar_id"`
	Name     string `json:"name"`
}

func ToAvatarDTO(a *domain.User) (dto AvatarDTO) {
	dto.AvatarId = a.AvatarId.AvatarId()
	dto.Name = a.Account.Account()
	return
}

func newUserDTO(u *domain.User) (dto UserDTO) {
	dto.Account = u.Account.Account()
	if u.AvatarId != nil {
		dto.AvatarId = u.AvatarId.AvatarId()
	}

	if u.Desc != nil {
		dto.Bio = u.Desc.MSDDesc()
	}

	dto.Email = u.Email.Email()

	dto.Password = u.PlatformPwd
	dto.Fullname = u.Fullname.MSDFullname()
	dto.CreatedAt = u.CreatedAt
	dto.UpdatedAt = u.UpdatedAt
	dto.Id = u.Id.Identity()

	return
}

type UpdateUserBasicInfoCmd struct {
	Desc            primitive.MSDDesc
	Email           primitive.Email
	AvatarId        primitive.AvatarId
	Fullname        primitive.MSDFullname
	descChanged     bool
	avatarChanged   bool
	emailChanged    bool
	fullNameChanged bool
}

func (cmd *UpdateUserBasicInfoCmd) toUser(u *domain.User) (changed bool) {
	if cmd.AvatarId != nil && cmd.AvatarId.AvatarId() != u.AvatarId.AvatarId() {
		u.AvatarId = cmd.AvatarId
		cmd.avatarChanged = true
	}

	if cmd.Desc != nil && cmd.Desc.MSDDesc() != u.Desc.MSDDesc() {
		u.Desc = cmd.Desc
		cmd.descChanged = true
	}

	if cmd.Email != nil && u.Email.Email() != cmd.Email.Email() {
		u.Email = cmd.Email
		cmd.emailChanged = true
	}

	if cmd.Fullname != nil && u.Fullname.MSDFullname() != cmd.Fullname.MSDFullname() {
		u.Fullname = cmd.Fullname
		cmd.fullNameChanged = true
	}

	changed = cmd.avatarChanged || cmd.descChanged || cmd.emailChanged || cmd.fullNameChanged

	return
}
