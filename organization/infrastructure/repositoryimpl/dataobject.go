package repositoryimpl

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
)

func toMemberDoc(o *domain.OrgMember) Member {
	do := Member{
		Username: o.Username.Account(),
		UserId:   o.UserId.Integer(),
		Orgname:  o.OrgName.Account(),
		OrgId:    o.OrgId.Integer(),
		Role:     o.Role,
		Version:  o.Version,
		Type:     o.Type,
	}

	do.ID = o.Id.Integer()
	return do
}

func toOrgMember(doc *Member) domain.OrgMember {
	return domain.OrgMember{
		Id:        primitive.CreateIdentity(doc.ID),
		OrgName:   primitive.CreateAccount(doc.Orgname),
		OrgId:     primitive.CreateIdentity(doc.OrgId),
		Role:      domain.OrgRole(doc.Role),
		Username:  primitive.CreateAccount(doc.Username),
		UserId:    primitive.CreateIdentity(doc.UserId),
		UpdatedAt: doc.CreatedAt.Unix(),
		CreatedAt: doc.CreatedAt.Unix(),
		Type:      doc.Type,
		Version:   doc.Version,
	}
}

func toApproveDoc(o *domain.Approve) Approve {
	do := Approve{
		Username: o.Username.Account(),
		UserId:   o.UserId.Integer(),
		Orgname:  o.OrgName.Account(),
		OrgId:    o.OrgId.Integer(),
		Role:     o.Role,
		Expire:   o.ExpireAt,
		Inviter:  o.Inviter.Account(),
		Status:   o.Status,
		By:       o.By,
		Msg:      o.Msg,
		Version:  o.Version,
		Type:     domain.InviteTypeInvite,
	}

	do.ID = o.Id.Integer()
	return do
}

func toRequestDoc(o *domain.MemberRequest) Approve {
	do := Approve{
		Username: o.Username.Account(),
		UserId:   o.UserId.Integer(),
		Orgname:  o.OrgName.Account(),
		OrgId:    o.OrgId.Integer(),
		Role:     o.Role,
		Status:   o.Status,
		By:       o.By,
		Msg:      o.Msg,
		Version:  o.Version,
		Type:     domain.InviteTypeRequest,
	}

	do.ID = o.Id.Integer()
	return do
}

func toApprove(doc *Approve) domain.Approve {
	return domain.Approve{
		Id:        primitive.CreateIdentity(doc.ID),
		Username:  primitive.CreateAccount(doc.Username),
		UserId:    primitive.CreateIdentity(doc.UserId),
		OrgName:   primitive.CreateAccount(doc.Orgname),
		OrgId:     primitive.CreateIdentity(doc.OrgId),
		Role:      domain.OrgRole(doc.Role),
		ExpireAt:  doc.Expire,
		Inviter:   primitive.CreateAccount(doc.Inviter),
		Version:   doc.Version,
		By:        doc.By,
		Status:    domain.ApproveStatus(doc.Status),
		Msg:       doc.Msg,
		CreatedAt: doc.CreatedAt.Unix(),
		UpdatedAt: doc.UpdatedAt.Unix(),
	}
}

func toMemberRequest(doc *Approve) domain.MemberRequest {
	return domain.MemberRequest{
		Id:        primitive.CreateIdentity(doc.ID),
		OrgName:   primitive.CreateAccount(doc.Orgname),
		OrgId:     primitive.CreateIdentity(doc.OrgId),
		Username:  primitive.CreateAccount(doc.Username),
		UserId:    primitive.CreateIdentity(doc.UserId),
		Role:      domain.OrgRole(doc.Role),
		Version:   doc.Version,
		By:        doc.By,
		Status:    domain.ApproveStatus(doc.Status),
		CreatedAt: doc.CreatedAt.Unix(),
		UpdatedAt: doc.UpdatedAt.Unix(),
		Msg:       doc.Msg,
	}
}
