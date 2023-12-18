package app

import (
	"fmt"

	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/repository"
)

type UserInfoDTO struct {
	Points int `json:"points"`

	UserDTO
}

type UserDTO struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Account string `json:"account"`

	Bio      string     `json:"bio"`
	AvatarId string     `json:"avatar_id"`
	Tokens   []TokenDTO `json:"tokens"`
	Password string     `json:"-"`
}

type TokenDTO struct {
	CreatedAt  int64  `json:"created_at"`
	Expire     int64  `json:"expired"`
	Account    string `json:"account"`
	Name       string `json:"name"`
	Permission string `json:"permission"`
	Token      string `json:"token"`
}

func newTokenDTO(t *domain.PlatformToken) (dto TokenDTO) {
	dto.CreatedAt = t.CreatedAt
	dto.Expire = t.Expire
	dto.Account = t.Account.Account()
	dto.Name = t.Name
	dto.Permission = string(t.Permission)
	dto.Token = t.Token

	return

}

func newUserDTO(u *domain.User) (dto UserDTO) {
	dto.Account = u.Account.Account()
	dto.AvatarId = u.AvatarId.AvatarId()
	dto.Bio = u.Bio.Bio()
	dto.Email = u.Email.Email()
	dto.Id = fmt.Sprint(u.PlatformId)

	dto.Tokens = make([]TokenDTO, 0, len(u.PlatformTokens))
	for name, t := range u.PlatformTokens {
		dto.Tokens = append(dto.Tokens, TokenDTO{
			CreatedAt:  t.CreatedAt,
			Expire:     t.Expire,
			Account:    u.Account.Account(),
			Name:       name,
			Permission: string(t.Permission),
		})
	}

	dto.Password = u.PlatformPwd

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

type FollowsListCmd struct {
	User domain.Account

	repository.FollowFindOption
}

type FollowsDTO struct {
	Total int         `json:"total"`
	Data  []FollowDTO `json:"data"`
}

type FollowDTO struct {
	Account    string `json:"account"`
	AvatarId   string `json:"avatar_id"`
	Bio        string `json:"bio"`
	IsFollower bool   `json:"is_follower"`
}
