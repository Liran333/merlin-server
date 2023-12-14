package repositoryimpl

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
)

type userInfo struct {
	DUser `bson:",inline"`

	Count int `bson:"count"`
}

func toUserDoc(u domain.User, doc *DUser) {
	*doc = DUser{
		Name:     u.Account.Account(),
		Email:    u.Email.Email(),
		AvatarId: u.AvatarId.AvatarId(),
	}

	if u.Bio != nil {
		doc.Bio = u.Bio.Bio()
	}
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
	u.Version = doc.Version

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
