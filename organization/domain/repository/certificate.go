/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
	orgprimitive "github.com/openmerlin/merlin-server/organization/domain/primitive"
)

// FindOption represents the options for finding an organization certificate.
type FindOption struct {
	Phone                   primitive.Phone
	OrgName                 primitive.Account
	CertificateOrgName      primitive.AccountFullname
	UnifiedSocialCreditCode orgprimitive.USCC
}

// Certificate represents the repository for organization certificates.
type Certificate interface {
	Save(domain.OrgCertificate) error
	Find(FindOption) (domain.OrgCertificate, error)
	DuplicateCheck(option FindOption) (domain.OrgCertificate, error)
	DeleteByOrgName(orgName primitive.Account) error
}
