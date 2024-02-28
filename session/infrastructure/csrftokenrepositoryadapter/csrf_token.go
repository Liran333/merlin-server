/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package csrftokenrepositoryadapter provides an adapter for the CSRF token repository.
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

// NewCSRFTokenAdapter creates a new instance of the CSRF token adapter with the given DAO.
func NewCSRFTokenAdapter(d dao) *csrfTokenAdapter {
	return &csrfTokenAdapter{
		dao: d,
	}
}

type csrfTokenAdapter struct {
	dao dao
}

// Add adds a CSRF token to the repository.
func (adapter *csrfTokenAdapter) Add(id primitive.RandomId, t *domain.CSRFToken) error {
	v := toCSRFTokenDO(t)

	// must pass *csrfTokenDO, because it implements the interface of MarshalBinary
	return adapter.dao.SetWithExpiry(id.RandomId(), &v, t.LifeTime())
}

// Find finds a CSRF token in the repository by its ID.
func (adapter *csrfTokenAdapter) Find(tid primitive.RandomId) (domain.CSRFToken, error) {
	var do csrfTokenDO

	// note the *csrfTokenDO implements interface of UnmarshalBinary
	if err := adapter.dao.Get(tid.RandomId(), &do); err != nil {
		if adapter.dao.IsKeyNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}

		return domain.CSRFToken{}, err
	}

	return do.toCSRFToken(), nil
}
