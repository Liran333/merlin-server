package csrftokenrepositoryadapter

import (
	"time"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/session/domain"
)

type dao interface {
	Get(key string, val interface{}) error
	SetWithExpiry(key string, val interface{}, expiry time.Duration) error
	IsKeyNotExists(err error) bool
}

func NewCSRFTokenAdapter(d dao) *csrfTokenAdapter {
	return &csrfTokenAdapter{
		dao: d,
	}
}

type csrfTokenAdapter struct {
	dao dao
}

func (adapter *csrfTokenAdapter) Add(t *domain.CSRFToken) error {
	v := toCSRFTokenDO(t)

	// must pass *csrfTokenDO, because it implements the interface of MarshalBinary
	return adapter.dao.SetWithExpiry(t.Id.String(), &v, t.LifeTime())
}

func (adapter *csrfTokenAdapter) Save(t *domain.CSRFToken) error {
	return adapter.Add(t)
}

func (adapter *csrfTokenAdapter) Find(tid primitive.UUID) (domain.CSRFToken, error) {
	var do csrfTokenDO

	// note the *csrfTokenDO implements interface of UnmarshalBinary
	if err := adapter.dao.Get(tid.String(), &do); err != nil {
		if adapter.dao.IsKeyNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}

		return domain.CSRFToken{}, err
	}

	return do.toCSRFToken(tid), nil
}
