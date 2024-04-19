package app

import (
	"errors"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepository "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/organization/domain/email"
	"github.com/openmerlin/merlin-server/organization/domain/repository"
)

type OrgCertificateService interface {
	Certificate(*OrgCertificateCmd) error
	GetCertification(orgName, actor primitive.Account) (OrgCertificateDTO, error)
	DuplicateCheck(cmd OrgCertificateDuplicateCheckCmd) (bool, error)
}

func NewOrgCertificateService(
	perm *permService,
	email email.Email,
	cert repository.Certificate,
) OrgCertificateService {
	return &orgCertificateService{
		perm:        perm,
		email:       email,
		certificate: cert,
	}
}

type orgCertificateService struct {
	perm        *permService
	email       email.Email
	certificate repository.Certificate
}

func (org *orgCertificateService) Certificate(cmd *OrgCertificateCmd) error {
	err := org.perm.Check(cmd.Actor, cmd.OrgName, primitive.ObjTypeOrg, primitive.ActionWrite)
	if err != nil {
		return err
	}

	option := repository.FindOption{
		Phone:                   cmd.Phone,
		CertificateOrgName:      cmd.CertificateOrgName,
		UnifiedSocialCreditCode: cmd.UnifiedSocialCreditCode,
	}

	_, err = org.certificate.DuplicateCheck(option)
	if err == nil {
		return errors.New("duplicate information")
	}
	if !commonrepository.IsErrorResourceNotExists(err) {
		return err
	}

	certificateData := cmd.OrgCertificate
	certificateData.SetProcessingStatus()

	if err = org.certificate.Save(certificateData); err != nil {
		return err
	}

	return org.email.Send(cmd.OrgCertificate, cmd.ImageOfCertificate)
}

func (org *orgCertificateService) GetCertification(orgName, actor primitive.Account) (OrgCertificateDTO, error) {
	cert, err := org.certificate.Find(repository.FindOption{OrgName: orgName})
	if err != nil {
		if commonrepository.IsErrorResourceNotExists(err) {
			err = nil
		}

		return OrgCertificateDTO{}, err
	}

	isAdmin := false
	if actor != nil {
		err = org.perm.Check(actor, orgName, primitive.ObjTypeOrg, primitive.ActionWrite)
		if err == nil {
			isAdmin = true
		}
	}

	dto := toCertificationDTO(cert)
	if !isAdmin {
		dto.Masked()
	}

	return dto, nil
}

func (org *orgCertificateService) DuplicateCheck(cmd OrgCertificateDuplicateCheckCmd) (bool, error) {
	_, err := org.certificate.DuplicateCheck(cmd)
	if err != nil {
		if commonrepository.IsErrorResourceNotExists(err) {
			return true, nil
		}
	}

	return false, err
}
