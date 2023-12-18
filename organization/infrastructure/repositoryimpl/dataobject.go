package repositoryimpl

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
	userdomain "github.com/openmerlin/merlin-server/user/domain"
)

func toOrgDoc(o domain.Organization) Organization {
	approves := make([]Approve, len(o.Approves))
	for i, approve := range o.Approves {
		approves[i] = ToApproveDoc(approve)
	}
	return Organization{
		Name:        o.Name.Account(),
		AvatarId:    o.AvatarId.AvatarId(),
		FullName:    o.FullName,
		Website:     o.Website,
		Description: o.Description,
		PlatformId:  o.PlatformId,
		Owner:       o.Owner.Account(),
		OwnerTeamId: o.OwnerTeamId,
		WriteTeamId: o.WriteTeamId,
		ReadTeamId:  o.ReadTeamId,
		Approves:    approves,
	}
}

func toOrganization(doc Organization, u *domain.Organization) (err error) {

	u.Name = primitive.CreateAccount(doc.Name)

	if u.AvatarId, err = userdomain.NewAvatarId(doc.AvatarId); err != nil {
		return
	}

	u.Id = doc.Id.Hex()
	u.PlatformId = doc.PlatformId
	u.Version = doc.Version
	u.FullName = doc.FullName
	u.Website = doc.Website
	u.Description = doc.Description
	u.OwnerTeamId = doc.OwnerTeamId
	u.ReadTeamId = doc.ReadTeamId
	u.WriteTeamId = doc.WriteTeamId
	u.Owner = primitive.CreateAccount(doc.Owner)

	u.Approves = make([]domain.Approve, len(doc.Approves))
	for i := range doc.Approves {
		u.Approves[i] = toApprove(&doc.Approves[i])
	}

	return
}

func toMemberDoc(o domain.OrgMember) Member {

	return Member{
		Username: o.Username,
		Orgname:  o.OrgName,
		Role:     string(o.Role),
	}
}

func toOrgMember(doc *Member) domain.OrgMember {
	return domain.OrgMember{
		Id:       doc.Id.Hex(),
		OrgName:  doc.Orgname,
		Role:     domain.OrgRole(doc.Role),
		Username: doc.Username,
		Version:  doc.Version,
	}
}

func ToApproveDoc(o domain.Approve) Approve {
	return Approve{
		Username: o.Username,
		Orgname:  o.OrgName,
		Role:     string(o.Role),
		Expire:   o.ExpireAt,
	}
}

func toApprove(doc *Approve) domain.Approve {
	return domain.Approve{
		Username: doc.Username,
		OrgName:  doc.Orgname,
		Role:     domain.OrgRole(doc.Role),
		ExpireAt: doc.Expire,
	}
}
