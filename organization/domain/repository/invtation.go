package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
)

type Approve interface {
	AddInvite(*domain.Approve) (domain.Approve, error)
	SaveInvite(*domain.Approve) (domain.Approve, error)
	AddRequest(*domain.MemberRequest) (domain.MemberRequest, error)
	SaveRequest(*domain.MemberRequest) (domain.MemberRequest, error)
	DeleteInviteAndReqByOrg(primitive.Account) error
	//DeleteRequestByOrg(primitive.Account) error
	ListInvitation(*domain.OrgInvitationListCmd) ([]domain.Approve, error)
	ListRequests(*domain.OrgMemberReqListCmd) ([]domain.MemberRequest, error)
}
