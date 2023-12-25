package repositoryimpl

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
	userdomain "github.com/openmerlin/merlin-server/user/domain"
)

func toOrgDoc(o domain.Organization) Organization {
	return Organization{
		Name:              o.Name.Account(),
		AvatarId:          o.AvatarId.AvatarId(),
		FullName:          o.FullName,
		Website:           o.Website,
		Description:       o.Description,
		PlatformId:        o.PlatformId,
		Owner:             o.Owner.Account(),
		OwnerTeamId:       o.OwnerTeamId,
		WriteTeamId:       o.WriteTeamId,
		ReadTeamId:        o.ReadTeamId,
		ContributorTeamId: o.ContributorTeamId,
		DefaultRole:       string(o.DefaultRole),
		AllowRequest:      o.AllowRequest,
		Type:              int(primitive.OrgType),
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
	u.ContributorTeamId = doc.ContributorTeamId
	u.OwnerTeamId = doc.OwnerTeamId
	u.ReadTeamId = doc.ReadTeamId
	u.WriteTeamId = doc.WriteTeamId
	u.Owner = primitive.CreateAccount(doc.Owner)
	u.Type = doc.Type
	u.DefaultRole = domain.OrgRole(doc.DefaultRole)
	u.AllowRequest = doc.AllowRequest

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
		Username:  o.Username,
		Orgname:   o.OrgName,
		Role:      string(o.Role),
		Expire:    o.ExpireAt,
		Inviter:   o.Inviter,
		Status:    string(o.Status),
		By:        o.By,
		Msg:       o.Msg,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func toApprove(doc *Approve) domain.Approve {
	return domain.Approve{
		Id:        doc.Id.Hex(),
		Username:  doc.Username,
		OrgName:   doc.Orgname,
		Role:      domain.OrgRole(doc.Role),
		ExpireAt:  doc.Expire,
		Inviter:   doc.Inviter,
		Version:   doc.Version,
		By:        doc.By,
		Status:    domain.ApproveStatus(doc.Status),
		Msg:       doc.Msg,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}

func toMemberRequest(doc *MemberRequest) domain.MemberRequest {
	return domain.MemberRequest{
		Id:        doc.Id.Hex(),
		OrgName:   doc.Orgname,
		Username:  doc.Username,
		Role:      domain.OrgRole(doc.Role),
		Version:   doc.Version,
		By:        doc.By,
		Status:    domain.ApproveStatus(doc.Status),
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
		Msg:       doc.Msg,
	}
}

func toMemberRequestDoc(o domain.MemberRequest) MemberRequest {
	return MemberRequest{
		Orgname:   o.OrgName,
		Username:  o.Username,
		Role:      string(o.Role),
		Status:    string(o.Status),
		By:        o.By,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		Msg:       o.Msg,
	}
}
