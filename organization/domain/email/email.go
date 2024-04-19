package email

import (
	"github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/organization/domain/primitive"
)

type Email interface {
	Send(domain.OrgCertificate, primitive.Image) error
}
