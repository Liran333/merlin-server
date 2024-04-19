package repositoryimpl

import (
	"github.com/openmerlin/merlin-server/common/domain/crypto"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/organization/domain"
	orgprimitive "github.com/openmerlin/merlin-server/organization/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain/repository"
)

func NewCertificateImpl(db postgresql.Impl, enc crypto.Encrypter) (*certificateRepoImpl, error) {
	certificateTableName = db.TableName()
	err := db.DB().AutoMigrate(&CertificateDO{})

	return &certificateRepoImpl{Impl: db, e: enc}, err
}

type certificateRepoImpl struct {
	postgresql.Impl
	e crypto.Encrypter
}

func (impl *certificateRepoImpl) Save(cert domain.OrgCertificate) error {
	do, err := toCertificateDo(cert, impl.e)
	if err != nil {
		return err
	}

	return impl.DB().Save(&do).Error
}

func (impl *certificateRepoImpl) Find(option repository.FindOption) (domain.OrgCertificate, error) {
	do := CertificateDO{}
	if option.OrgName != nil {
		do.OrgName = option.OrgName.Account()
	}

	if option.CertificateOrgName != nil {
		do.CertificateOrgName = option.CertificateOrgName.AccountFullname()
	}

	if option.UnifiedSocialCreditCode != nil {
		do.USCC = option.UnifiedSocialCreditCode.USCC()
	}

	if err := impl.GetRecord(&do, &do); err != nil {
		return domain.OrgCertificate{}, err
	}

	return do.toCertificate(impl.e)
}

func (impl *certificateRepoImpl) DuplicateCheck(option repository.FindOption) (domain.OrgCertificate, error) {
	var do CertificateDO

	queryOr := impl.DB().Order(fieldID)
	if option.OrgName != nil {
		queryOr.Or(impl.EqualQuery(fieldOrg), option.OrgName.Account())
	}

	if option.CertificateOrgName != nil {
		queryOr.Or(impl.EqualQuery(fieldCertOrgName), option.CertificateOrgName.AccountFullname())
	}

	if option.UnifiedSocialCreditCode != nil {
		queryOr.Or(impl.EqualQuery(fieldUSCC), option.UnifiedSocialCreditCode.USCC())
	}

	if option.Phone != nil {
		queryOr.Or(impl.EqualQuery(fieldPhone), option.Phone.PhoneNumber())
	}

	query := impl.DB().
		Where(impl.EqualQuery(fieldStatus), orgprimitive.NewPassedStatus().CertificateStatus()).
		Where(queryOr)

	if err := impl.GetRecord(query, &do); err != nil {
		return domain.OrgCertificate{}, err
	}

	return do.toCertificate(impl.e)
}

func (impl *certificateRepoImpl) DeleteByOrgName(orgName primitive.Account) error {
	return impl.DB().Where(impl.EqualQuery(fieldOrg), orgName.Account()).Delete(&CertificateDO{}).Error
}
