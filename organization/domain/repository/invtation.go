package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
)

type MemberRequest interface {
	Save(*domain.MemberRequest) (domain.MemberRequest, error)
	DeleteByOrg(primitive.Account) error
	ListInvitation(*domain.OrgMemberReqListCmd) ([]domain.MemberRequest, error)
	//Search(*UserSearchOption) (UserSearchResult, error)
}

type Approve interface {
	Save(*domain.Approve) (domain.Approve, error)
	DeleteByOrg(primitive.Account) error
	ListInvitation(*domain.OrgInvitationListCmd) ([]domain.Approve, error)
	//Search(*UserSearchOption) (UserSearchResult, error)
}
