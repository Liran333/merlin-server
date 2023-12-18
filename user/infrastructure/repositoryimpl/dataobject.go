package repositoryimpl

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
)

type userInfo struct {
	DUser `bson:",inline"`

	Count int `bson:"count"`
}

func toUserDoc(u domain.User, doc *DUser) error {

	*doc = DUser{
		Name:        u.Account.Account(),
		Email:       u.Email.Email(),
		AvatarId:    u.AvatarId.AvatarId(),
		PlatformId:  u.PlatformId,
		PlatformPwd: u.PlatformPwd,
	}

	doc.PlatformTokens = make([]DToken, 0, len(u.PlatformTokens))
	for name, t := range u.PlatformTokens {

		doc.PlatformTokens = append(doc.PlatformTokens, DToken{
			Name:       t.Name,
			Expire:     t.Expire,
			CreatedAt:  t.CreatedAt,
			Account:    name,
			Permission: string(t.Permission),
		})

	}

	if u.Bio != nil {
		doc.Bio = u.Bio.Bio()
	}
	return nil
}

func toUser(doc DUser, u *domain.User) (err error) {

	if u.Email, err = domain.NewEmail(doc.Email); err != nil {
		return
	}

	if u.Account, err = primitive.NewAccount(doc.Name); err != nil {
		return
	}

	if u.Bio, err = domain.NewBio(doc.Bio); err != nil {
		return
	}

	if u.AvatarId, err = domain.NewAvatarId(doc.AvatarId); err != nil {
		return
	}

	u.Id = doc.Id.Hex()
	u.PlatformPwd = doc.PlatformPwd
	u.PlatformId = doc.PlatformId
	u.Version = doc.Version

	u.PlatformTokens = make(map[string]domain.PlatformToken)
	for _, t := range doc.PlatformTokens {
		u.PlatformTokens[t.Name] = domain.PlatformToken{
			CreatedAt:  t.CreatedAt,
			Expire:     t.Expire,
			Name:       t.Name,
			Permission: domain.TokenPerm(t.Permission),
			Account:    u.Account,
		}
	}

	return
}

func toUserInfo(doc DUser, info *domain.UserInfo) (err error) {

	if info.Account, err = primitive.NewAccount(doc.Name); err != nil {
		return
	}

	if info.AvatarId, err = domain.NewAvatarId(doc.AvatarId); err != nil {
		return
	}

	return
}

func toFollowerUserInfo(doc DUser, info *domain.FollowerUserInfo) (err error) {

	if info.Account, err = primitive.NewAccount(doc.Name); err != nil {
		return
	}

	if info.AvatarId, err = domain.NewAvatarId(doc.AvatarId); err != nil {
		return
	}

	if info.Bio, err = domain.NewBio(doc.Bio); err != nil {
		return
	}

	return
}
