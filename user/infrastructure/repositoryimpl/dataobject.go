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
		Fullname:    u.Fullname,
	}

	if u.Bio != nil {
		doc.Bio = u.Bio.Bio()
	}
	return nil
}

func toTokenDoc(t domain.PlatformToken, doc *DToken) error {

	*doc = DToken{
		Name:       t.Name,
		Expire:     t.Expire,
		CreatedAt:  t.CreatedAt,
		Account:    t.Account.Account(),
		Permission: string(t.Permission),
		Salt:       t.Salt,
		Token:      t.Token,
		LastEight:  t.LastEight,
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
	u.Fullname = doc.Fullname

	return
}

func toToken(doc DToken, t *domain.PlatformToken) {
	t.Name = doc.Name
	t.Account = primitive.CreateAccount(doc.Account)
	t.CreatedAt = doc.CreatedAt
	t.Expire = doc.Expire
	t.Salt = doc.Salt
	t.Token = doc.Token
	t.LastEight = doc.LastEight
	t.Permission = domain.TokenPerm(doc.Permission)
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
